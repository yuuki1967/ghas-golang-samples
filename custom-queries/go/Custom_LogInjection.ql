/**
 * @name Log entries created from user input
 * @description Building log entries from user-controlled sources is vulnerable to
 *              insertion of forged log entries by a malicious user.
 * @kind path-problem
 * @problem.severity error
 * @security-severity 7.8
 * @precision high
 * @id go/customs-log-injection
 * @tags security
 *       external/cwe/cwe-117
 */

import go
import DataFlow::PathGraph

/**
 * Provides a taint-tracking configuration for reasoning about
 * log injection vulnerabilities.
 */
module LogrusInjection {
  /**
   * A data flow source for log injection vulnerabilities.
   */
  abstract class Source extends DataFlow::Node { }

  /**
   * A data flow sink for log injection vulnerabilities.
   */
  abstract class Sink extends DataFlow::Node { }

  /**
   * A sanitizer for log injection vulnerabilities.
   */
  abstract class Sanitizer extends DataFlow::Node { }

  /**
   * DEPRECATED: Use `Sanitizer` instead.
   *
   * A sanitizer guard for log injection vulnerabilities.
   */
  abstract deprecated class SanitizerGuard extends DataFlow::BarrierGuard { }

  /** A source of untrusted data, considered as a taint source for log injection. */
  class UntrustedFlowAsSource extends Source {
    UntrustedFlowAsSource() { this instanceof UntrustedFlowSource }
  }

  /** An argument to a logging mechanism. */
  class LoggerSink extends Sink {
    LoggerSink() { this = any(LoggerCall log).getAMessageComponent() }
  }

  /**
   * A call to `strings.Replace` or `strings.ReplaceAll`, considered as a sanitizer
   * for log injection.
  class ReplaceSanitizer extends Sanitizer {
    ReplaceSanitizer() {
      exists(string name, DataFlow::CallNode call |
        this = call and
        call.getTarget().hasQualifiedName("strings", name) and
        call.getArgument(1).getStringValue().matches("%" + ["\r", "\n"] + "%")
      |
        name = "Replace" and
        call.getArgument(3).getNumericValue() < 0
        or
        name = "ReplaceAll"
      )
    }
  }

  /**
   * An argument that is formatted using the `%q` directive, considered as a sanitizer
   * for log injection.
   *
   * This formatting directive replaces newline characters with escape sequences.
  */
  private class SafeFormatArgumentSanitizer extends Sanitizer {
    SafeFormatArgumentSanitizer() {
      exists(StringOps::Formatting::StringFormatCall call, string safeDirective |
        this = call.getOperand(_, safeDirective) and
        // Mark "%q" formats as safe, but not "%#q", which would preserve newline characters.
        safeDirective.regexpMatch("%[^%#]*q")
      ) or
      exists(CallExpr cexpr | cexpr.getAnArgument().getAChild().getChild(1).getChild(0).toString() = "DisableHTMLEscape" and 
        cexpr.getAnArgument().getAChild().getChild(1).getChild(1).toString() = "false")
    }
  }
  /**
   * A taint-tracking configuration for reasoning about log injection vulnerabilities.
   */
  class Configuration extends TaintTracking::Configuration {
    Configuration() { this = "LogInjection" }

    override predicate isSource(DataFlow::Node source) { source instanceof Source }

    override predicate isSink(DataFlow::Node sink) { sink instanceof Sink }

    override predicate isSanitizer(DataFlow::Node sanitizer) { sanitizer instanceof Sanitizer }

    deprecated override predicate isSanitizerGuard(DataFlow::BarrierGuard guard) {
      guard instanceof SanitizerGuard
    }
  }
}

from LogrusInjection::Configuration c, DataFlow::PathNode source, DataFlow::PathNode sink
where c.hasFlowPath(source, sink)
select sink.getNode(), source, sink, "This log entry depends on a $@.", source.getNode(),
  "user-provided value"