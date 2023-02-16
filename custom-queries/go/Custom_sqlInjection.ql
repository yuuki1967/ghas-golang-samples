/**
 * @name Database query built from user-controlled sources
 * @description Building a database query from user-controlled sources is vulnerable to insertion of
 *              malicious code by the user.
 * @kind path-problem
 * @problem.severity error
 * @security-severity 8.8
 * @precision high
 * @id go/custom-sql-injection
 * @tags security
 *       external/cwe/cwe-089
 */

 import go
// import semmle.go.security.MongoSqlInjection
 import DataFlow::PathGraph
/**
 * Provides extension points for customizing the taint tracking configuration for reasoning about
 * SQL-injection vulnerabilities.
 */
module MongoSqlInjection {
  /**
   * A data flow source for SQL-injection vulnerabilities.
   */
  abstract class Source extends DataFlow::Node { }

  /**
   * A data flow sink for SQL-injection vulnerabilities.
   */
  abstract class Sink extends DataFlow::Node { }

  /**
   * A sanitizer for SQL-injection vulnerabilities.
   */
  abstract class Sanitizer extends DataFlow::Node { }

  /**
   * DEPRECATED: Use `Sanitizer` instead.
   *
   * A sanitizer guard for SQL-injection vulnerabilities.
   */
  abstract deprecated class SanitizerGuard extends DataFlow::BarrierGuard { 
  }

  /** A source of untrusted data, considered as a taint source for SQL injection. */
  class UntrustedFlowAsSource extends Source instanceof UntrustedFlowSource { }

  /** An SQL string, considered as a taint sink for SQL injection. */
  class SqlQueryAsSink extends Sink instanceof SQL::QueryString { }

  /** A NoSql query, considered as a taint sink for SQL injection. */
  class NoSqlQueryAsSink extends Sink instanceof NoSql::Query { }
  private class MongoBsonSanitizer extends Sanitizer {
    MongoBsonSanitizer() {
      exists(DataFlow::Node node, SliceLit slit,SelectorExpr sexpr |
        this = node and slit = node.asExpr() and 
        slit.getAChild().getAChild().toString().matches("bson"))
    }
  }
  /**
   * A taint-tracking configuration for reasoning about SQL-injection vulnerabilities.
   */
  class Configuration extends TaintTracking::Configuration {
    Configuration() { this = "MongoSqlInjection" }

    override predicate isSource(DataFlow::Node source) { source instanceof Source }

    override predicate isSink(DataFlow::Node sink) { sink instanceof Sink }

    override predicate isAdditionalTaintStep(DataFlow::Node pred, DataFlow::Node succ) {
      NoSql::isAdditionalMongoTaintStep(pred, succ)
    }

    override predicate isSanitizer(DataFlow::Node node) {
      super.isSanitizer(node) or
      node instanceof Sanitizer
    }

    deprecated override predicate isSanitizerGuard(DataFlow::BarrierGuard guard) {
      guard instanceof SanitizerGuard
    }
  }
}
from MongoSqlInjection::Configuration cfg, DataFlow::PathNode source, DataFlow::PathNode sink
 where cfg.hasFlowPath(source, sink)
 select sink.getNode(), source, sink, "This query depends on a $@.", source.getNode(),
   "user-provided value"