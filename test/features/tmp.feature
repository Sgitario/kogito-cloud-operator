@springboot
Feature: Deploy spring boot service

  Background:
    Given Namespace is created

  @security
  Scenario: Deploy process-springboot-example service to complete user tasks
    Given Kogito Operator is deployed with Keycloak operator
    When Install Kogito Infra "Keycloak"

    Given Clone Kogito examples into local directory
    And Local example service "process-usertasks-with-security-oidc-springboot" is built by Maven
    When Create springboot service "process-usertasks-with-security-oidc-springboot" with configuration:
      | runtime-env | keycloak.auth-server-url | http://keycloak:8281/auth  |
      | runtime-env | keycloak.auth-server-url | http://keycloak:8281/auth  |
    And BuildConfig "process-usertasks-with-security-oidc-springboot-binary" is created after 1 minutes
    And Start build with name "process-usertasks-with-security-oidc-springboot-binary" from local example service path "process-usertasks-with-security-oidc-springboot/target"

    # And Deploy springboot example service "process-usertasks-with-security-oidc-springboot" with configuration:
    #  | runtime-env | keycloak.auth-server-url | http://keycloak:8281/auth  |
    And Kogito application "process-usertasks-with-security-oidc-springboot" has 1 pods running within 10 minutes

    When Start "approvals" process on service "process-usertasks-with-security-oidc-springboot" within 20 minutes with body:
      """json
      {
      "traveller" : {
        "firstName" : "John",
        "lastName" : "Doe",
        "email" : "jon.doe@example.com",
        "nationality" : "American",
        "address" : {
          "street" : "main street",
          "city" : "Boston",
          "zipCode" : "10005",
          "country" : "US" }
        }
      }
      """
