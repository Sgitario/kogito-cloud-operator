@smoke
Feature: Deploy Kogito Runtime

  Background:
    Given Namespace is created

  Scenario Outline: Deploy Simplest Scenario using Kogito Runtime
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

  @persistence
  Scenario Outline: Deploy Scenario with persistence using Kogito Runtime
    Given Kogito Operator is deployed with Infinispan operator

    When Deploy <runtime> example runtime service "quay.io/jcarvaja/<image>:8.0.0" with configuration:
      | config | persistence | enabled |
    And Kogito application "<image>" has 1 pods running within 10 minutes
    And Start "orders" process on service "<image>" with body:
      """json
      {
        "approver" : "john", 
        "order" : {
          "orderNumber" : "12345", 
          "shipped" : false
        }
      }
      """
    
    Then Service "<image>" contains 1 instances of process with name "orders"

    When Scale Kogito Runtime "<image>" to 0 pods within 2 minutes
    And Scale Kogito Runtime "<image>" to 1 pods within 2 minutes

    Then Service "<image>" contains 1 instances of process with name "orders" within 2 minutes

    @springboot
    Examples:
      | runtime    | image                                      |
      | springboot | process-springboot-example-persistence     |

    @quarkus
    Examples:
      | runtime    | image                                      |
      | quarkus    | process-quarkus-example-jvm-persistence    |

    @quarkus
    @native
    Examples:
      | runtime    | image                                      |
      | quarkus    | process-quarkus-example-native-persistence |