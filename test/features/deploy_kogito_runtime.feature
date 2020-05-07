@smoke
Feature: Deploy Kogito Runtime

  Background:
    Given Namespace is created

  Scenario Outline: Deploy Simplest Scenario with Kogito Runtime
    Given Kogito Operator is deployed

    When Deploy <runtime> example runtime service "quay.io/jcarvaja/<image>:8.0.0" with configuration:
      | config | persistence | disabled |
    
    Then Kogito application "<image>" has 1 pods running within 10 minutes
    And Service "<image>" with process name "orders" is available within 2 minutes

    @springboot
    Examples:
      | runtime    | image                          |
      | springboot | process-springboot-example     |

    @quarkus
    Examples:
      | runtime    | image                          |
      | quarkus    | process-quarkus-example-jvm    |

    @quarkus
    @native
    Examples:
      | runtime    | image                          |
      | quarkus    | process-quarkus-example-native |
