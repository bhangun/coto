Coto Output
Generated: 2026-01-17 09:24:42
Files: 17 | Directories: 6 | Total Size: 41.8 KB


================================================================================
tech/kayys/silat/sdk/client/CreateRunBuilder.java
Size: 2.7 KB | Modified: 2026-01-06 16:44:05
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.RunResponse;
import tech.kayys.silat.model.CreateRunRequest;

import java.util.HashMap;
import java.util.Map;

/**
 * Builder for creating workflow runs
 */
public class CreateRunBuilder {

    private final WorkflowRunClient client;
    private final String workflowDefinitionId;
    private final Map<String, Object> inputs = new HashMap<>();
    private final Map<String, String> labels = new HashMap<>();
    private String workflowVersion = "1.0.0";
    private String correlationId;
    private boolean autoStart = false;

    CreateRunBuilder(WorkflowRunClient client, String workflowDefinitionId) {
        this.client = client;
        this.workflowDefinitionId = workflowDefinitionId;
    }

    public CreateRunBuilder version(String version) {
        this.workflowVersion = version;
        return this;
    }

    public CreateRunBuilder input(String key, Object value) {
        inputs.put(key, value);
        return this;
    }

    public CreateRunBuilder inputs(Map<String, Object> inputs) {
        this.inputs.putAll(inputs);
        return this;
    }

    public CreateRunBuilder correlationId(String correlationId) {
        this.correlationId = correlationId;
        return this;
    }

    public CreateRunBuilder autoStart(boolean autoStart) {
        this.autoStart = autoStart;
        return this;
    }

    public CreateRunBuilder label(String key, String value) {
        if (key == null || key.trim().isEmpty()) {
            throw new IllegalArgumentException("Label key cannot be null or empty");
        }
        if (value == null) {
            throw new IllegalArgumentException("Label value cannot be null");
        }
        this.labels.put(key, value);
        return this;
    }

    public CreateRunBuilder labels(Map<String, String> labels) {
        if (labels != null) {
            for (Map.Entry<String, String> entry : labels.entrySet()) {
                label(entry.getKey(), entry.getValue());
            }
        }
        return this;
    }

    /**
     * Get the labels map for validation or debugging purposes
     */
    public Map<String, String> getLabels() {
        return new HashMap<>(labels);
    }

    /**
     * Execute and return the created run
     */
    public Uni<RunResponse> execute() {
        CreateRunRequest request = new CreateRunRequest(
                workflowDefinitionId,
                workflowVersion,
                inputs,
                correlationId,
                autoStart);
        return client.createRun(request);
    }

    /**
     * Execute and immediately start the run
     */
    public Uni<RunResponse> executeAndStart() {
        this.autoStart = true;
        return execute();
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/GrpcWorkflowDefinitionClient.java
Size: 1.5 KB | Modified: 2026-01-06 16:47:42
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.WorkflowDefinition;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * gRPC-based workflow definition client
 */
class GrpcWorkflowDefinitionClient implements WorkflowDefinitionClient {

    private final SilatClientConfig config;
    private final AtomicBoolean closed = new AtomicBoolean(false);

    GrpcWorkflowDefinitionClient(SilatClientConfig config) {
        this.config = config;
    }

    /**
     * Get the client configuration
     */
    public SilatClientConfig config() {
        return config;
    }

    // Implement using gRPC stubs...

    @Override
    public Uni<WorkflowDefinition> createDefinition(WorkflowDefinition request) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<WorkflowDefinition> getDefinition(String definitionId) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<List<WorkflowDefinition>> listDefinitions(boolean activeOnly) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<Void> deleteDefinition(String definitionId) {
        checkClosed();
        return null;
    }

    @Override
    public void close() {
        if (closed.compareAndSet(false, true)) {
            // Close gRPC resources if needed
        }
    }

    private void checkClosed() {
        if (closed.get()) {
            throw new IllegalStateException("Client is closed");
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/GrpcWorkflowRunClient.java
Size: 2.5 KB | Modified: 2026-01-06 16:47:43
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.RunResponse;
import tech.kayys.silat.model.CreateRunRequest;
import tech.kayys.silat.execution.ExecutionHistory;

import java.util.Map;
import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * gRPC-based workflow run client
 */
class GrpcWorkflowRunClient implements WorkflowRunClient {

    private final SilatClientConfig config;
    private final AtomicBoolean closed = new AtomicBoolean(false);
    // gRPC stub would be injected here

    GrpcWorkflowRunClient(SilatClientConfig config) {
        this.config = config;
    }

    /**
     * Get the client configuration
     */
    public SilatClientConfig config() {
        return config;
    }

    // Implement using gRPC stubs...

    @Override
    public Uni<RunResponse> createRun(CreateRunRequest request) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<RunResponse> getRun(String runId) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<RunResponse> startRun(String runId) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<RunResponse> suspendRun(String runId, String reason, String waitingOnNodeId) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<RunResponse> resumeRun(String runId, Map<String, Object> resumeData, String humanTaskId) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<Void> cancelRun(String runId, String reason) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<Void> signal(String runId, String signalName, String targetNodeId, Map<String, Object> payload) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<ExecutionHistory> getExecutionHistory(String runId) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<List<RunResponse>> queryRuns(String workflowId, String status, int page, int size) {
        checkClosed();
        return null;
    }

    @Override
    public Uni<Long> getActiveRunsCount() {
        checkClosed();
        return null;
    }

    @Override
    public void close() {
        if (closed.compareAndSet(false, true)) {
            // Close gRPC resources if needed
        }
    }

    private void checkClosed() {
        if (closed.get()) {
            throw new IllegalStateException("Client is closed");
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/QueryRunsBuilder.java
Size: 1001 B | Modified: 2026-01-03 10:41:05
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.RunResponse;
import java.util.List;

/**
 * Builder for querying runs
 */
public class QueryRunsBuilder {

    private final WorkflowRunClient client;
    private String workflowId;
    private String status;
    private int page = 0;
    private int size = 20;

    QueryRunsBuilder(WorkflowRunClient client) {
        this.client = client;
    }

    public QueryRunsBuilder workflowId(String workflowId) {
        this.workflowId = workflowId;
        return this;
    }

    public QueryRunsBuilder status(String status) {
        this.status = status;
        return this;
    }

    public QueryRunsBuilder page(int page) {
        this.page = page;
        return this;
    }

    public QueryRunsBuilder size(int size) {
        this.size = size;
        return this;
    }

    public Uni<List<RunResponse>> execute() {
        return client.queryRuns(workflowId, status, page, size);
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/RestWorkflowDefinitionClient.java
Size: 8.5 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import io.vertx.ext.web.client.WebClientOptions;
import io.vertx.mutiny.core.Vertx;
import io.vertx.mutiny.ext.web.client.WebClient;
import io.vertx.core.json.JsonObject;
import tech.kayys.silat.model.WorkflowDefinition;

import java.util.List;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * REST-based workflow definition client
 */
public class RestWorkflowDefinitionClient implements WorkflowDefinitionClient {

    private final SilatClientConfig config;
    private final Vertx vertx;
    private final WebClient webClient;
    private final AtomicBoolean closed = new AtomicBoolean(false);

    private static final com.fasterxml.jackson.databind.ObjectMapper mapper = new com.fasterxml.jackson.databind.ObjectMapper()
            .registerModule(new com.fasterxml.jackson.datatype.jsr310.JavaTimeModule())
            .registerModule(new com.fasterxml.jackson.module.paramnames.ParameterNamesModule())
            .configure(com.fasterxml.jackson.databind.DeserializationFeature.FAIL_ON_UNKNOWN_PROPERTIES, false);

    RestWorkflowDefinitionClient(SilatClientConfig config, Vertx vertx) {
        this.config = config;
        this.vertx = vertx;

        System.out.println("RestWorkflowDefinitionClient initialized with endpoint: '" + config.endpoint() + "'");
        System.out.println("Host: " + getHostFromEndpoint(config.endpoint()));
        System.out.println("Port: " + getPortFromEndpoint(config.endpoint()));

        // Use proper configuration
        WebClientOptions options = new WebClientOptions()
                .setDefaultHost(getHostFromEndpoint(config.endpoint()))
                .setDefaultPort(getPortFromEndpoint(config.endpoint()))
                .setSsl(config.endpoint().toLowerCase().startsWith("https"))
                .setConnectTimeout((int) config.timeout().toMillis())
                .setIdleTimeout((int) config.timeout().getSeconds());

        this.webClient = WebClient.create(vertx, options);
    }

    @Override
    public Uni<WorkflowDefinition> createDefinition(WorkflowDefinition request) {
        if (closed.get()) {
            return Uni.createFrom().failure(new IllegalStateException("Client is closed"));
        }

        tech.kayys.silat.dto.CreateWorkflowDefinitionRequest dto = tech.kayys.silat.dto.WorkflowDefinitionMapper
                .toCreateRequest(request);
        JsonObject requestBody = JsonObject.mapFrom(dto);

        return applyAuthHeaders(webClient
                .post(getPath("/api/v1/workflow-definitions"))
                .putHeader("Content-Type", "application/json")
                .putHeader("X-Tenant-ID", config.tenantId()))
                .sendJson(requestBody)
                .onItem().transform(response -> {
                    if (response.statusCode() == 200 || response.statusCode() == 201) {
                        String body = response.bodyAsString();
                        try {
                            return mapper.readValue(body, WorkflowDefinition.class);
                        } catch (Exception e) {
                            throw new RuntimeException("Failed to deserialize workflow definition: " + e.getMessage(),
                                    e);
                        }
                    }
                    throw new RuntimeException("Failed to create workflow definition: [" + response.statusCode() + "] "
                            + response.statusMessage() + " - " + response.bodyAsString());
                })
                .onFailure().transform(
                        msg -> new RuntimeException("Failed to create workflow definition: " + msg.getMessage(), msg));
    }

    @Override
    public Uni<WorkflowDefinition> getDefinition(String definitionId) {
        return applyAuthHeaders(webClient
                .get(getPath("/api/v1/workflow-definitions/" + definitionId))
                .putHeader("Accept", "application/json")
                .putHeader("X-Tenant-ID", config.tenantId()))
                .send()
                .onItem().transform(response -> {
                    try {
                        return mapper.readValue(response.bodyAsString(), WorkflowDefinition.class);
                    } catch (Exception e) {
                        throw new RuntimeException("Failed to deserialize workflow definition: " + e.getMessage(), e);
                    }
                })
                .onFailure().recoverWithUni(failure -> Uni.createFrom().failure(
                        new RuntimeException("Failed to get workflow definition: " + failure.getMessage(), failure)));
    }

    @Override
    public Uni<List<WorkflowDefinition>> listDefinitions(boolean activeOnly) {
        String query = activeOnly ? "?activeOnly=true" : "";
        return applyAuthHeaders(webClient
                .get(getPath("/api/v1/workflow-definitions" + query))
                .putHeader("Accept", "application/json")
                .putHeader("X-Tenant-ID", config.tenantId()))
                .send()
                .onItem().transform(response -> {
                    try {
                        return mapper.readValue(response.bodyAsString(),
                                new com.fasterxml.jackson.core.type.TypeReference<List<WorkflowDefinition>>() {
                                });
                    } catch (Exception e) {
                        throw new RuntimeException("Failed to deserialize workflow definitions: " + e.getMessage(), e);
                    }
                })
                .onFailure().recoverWithUni(failure -> Uni.createFrom().failure(
                        new RuntimeException("Failed to list workflow definitions: " + failure.getMessage(), failure)));
    }

    @Override
    public Uni<Void> deleteDefinition(String definitionId) {
        return applyAuthHeaders(webClient
                .delete(getPath("/api/v1/workflow-definitions/" + definitionId))
                .putHeader("X-Tenant-ID", config.tenantId()))
                .send()
                .onItem().transformToUni(response -> Uni.createFrom().voidItem())
                .onFailure().recoverWithUni(failure -> Uni.createFrom().failure(
                        new RuntimeException("Failed to delete workflow definition: " + failure.getMessage(),
                                failure)));
    }

    /**
     * Apply authentication headers based on configuration
     */
    private <T> io.vertx.mutiny.ext.web.client.HttpRequest<T> applyAuthHeaders(
            io.vertx.mutiny.ext.web.client.HttpRequest<T> request) {
        /*
         * if (config.apiKey() != null && !config.apiKey().trim().isEmpty()) {
         * request.putHeader("Authorization", "Bearer " + config.apiKey());
         * }
         */
        // Add any additional headers from config
        config.headers().forEach(request::putHeader);
        return request;
    }

    /**
     * Extract host from endpoint URL
     */
    private String getHostFromEndpoint(String endpoint) {
        if (endpoint.startsWith("http")) {
            return java.net.URI.create(endpoint).getHost();
        }
        // For host:port format
        int colonIndex = endpoint.indexOf(':');
        if (colonIndex != -1) {
            return endpoint.substring(0, colonIndex);
        }
        return endpoint;
    }

    /**
     * Extract port from endpoint URL
     */
    private int getPortFromEndpoint(String endpoint) {
        if (endpoint.startsWith("http")) {
            java.net.URI uri = java.net.URI.create(endpoint);
            int port = uri.getPort();
            if (port == -1) {
                return uri.getScheme().equals("https") ? 443 : 80;
            }
            return port;
        }
        // For host:port format
        int colonIndex = endpoint.indexOf(':');
        if (colonIndex != -1) {
            return Integer.parseInt(endpoint.substring(colonIndex + 1));
        }
        // Default to 80 for REST
        return 80;
    }

    /**
     * Get the API path, handling both absolute and relative endpoints
     */
    private String getPath(String path) {
        if (config.endpoint().startsWith("http")) {
            // If endpoint is a full URL, just return the path
            return path;
        } else {
            // If endpoint is host:port, prepend with "/"
            return path;
        }
    }

    /**
     * Close the client and release resources
     */
    public void close() {
        if (closed.compareAndSet(false, true)) {
            if (webClient != null) {
                webClient.close();
            }
            if (vertx != null) {
                vertx.close();
            }
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/RestWorkflowRunClient.java
Size: 6.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.RunResponse;
import tech.kayys.silat.model.CreateRunRequest;
import tech.kayys.silat.execution.ExecutionHistory;
import java.util.Map;
import java.util.List;

/**
 * REST-based workflow run client
 */
public class RestWorkflowRunClient implements WorkflowRunClient {

    private final SilatClientConfig config;
    private final io.vertx.mutiny.core.Vertx vertx;
    private final io.vertx.mutiny.ext.web.client.WebClient webClient;

    RestWorkflowRunClient(SilatClientConfig config, io.vertx.mutiny.core.Vertx vertx) {
        this.config = config;
        this.vertx = vertx;

        System.out.println("RestWorkflowRunClient initialized with endpoint: '" + config.endpoint() + "'");
        System.out.println("Host: " + getHostFromEndpoint(config.endpoint()));
        System.out.println("Port: " + getPortFromEndpoint(config.endpoint()));

        io.vertx.ext.web.client.WebClientOptions options = new io.vertx.ext.web.client.WebClientOptions()
                .setDefaultHost(getHostFromEndpoint(config.endpoint()))
                .setDefaultPort(getPortFromEndpoint(config.endpoint()))
                .setSsl(config.endpoint().toLowerCase().startsWith("https"));

        this.webClient = io.vertx.mutiny.ext.web.client.WebClient.create(vertx, options);
    }

    private String getHostFromEndpoint(String endpoint) {
        if (endpoint.startsWith("http")) {
            return java.net.URI.create(endpoint).getHost();
        }
        int colonIndex = endpoint.indexOf(':');
        if (colonIndex != -1) {
            return endpoint.substring(0, colonIndex);
        }
        return endpoint;
    }

    private int getPortFromEndpoint(String endpoint) {
        if (endpoint.startsWith("http")) {
            java.net.URI uri = java.net.URI.create(endpoint);
            int port = uri.getPort();
            if (port == -1) {
                return uri.getScheme().equals("https") ? 443 : 80;
            }
            return port;
        }
        int colonIndex = endpoint.indexOf(':');
        if (colonIndex != -1) {
            return Integer.parseInt(endpoint.substring(colonIndex + 1));
        }
        return 80;
    }

    @Override
    public Uni<RunResponse> createRun(CreateRunRequest request) {
        return webClient.post("/api/v1/workflow-runs")
                .putHeader("X-Tenant-ID", config.tenantId())
                // .putHeader("Authorization", "Bearer " + config.apiKey())
                .sendJson(request)
                .map(response -> {
                    System.out.println("RestWorkflowRunClient: createRun response status: " + response.statusCode());
                    System.out.println("RestWorkflowRunClient: createRun response body: " + response.bodyAsString());

                    io.vertx.core.json.JsonObject json = response.bodyAsJsonObject();
                    if (json == null) {
                        System.out.println("RestWorkflowRunClient: JSON is null!");
                        return null;
                    }

                    Object idObj = json.getValue("id");
                    String runId = (idObj instanceof io.vertx.core.json.JsonObject)
                            ? ((io.vertx.core.json.JsonObject) idObj).getString("value")
                            : (String) idObj;

                    Object defIdObj = json.getValue("definitionId");
                    String workflowId = (defIdObj instanceof io.vertx.core.json.JsonObject)
                            ? ((io.vertx.core.json.JsonObject) defIdObj).getString("value")
                            : (json.getString("workflowId") != null ? json.getString("workflowId") : (String) defIdObj);

                    RunResponse runResponse = RunResponse.builder()
                            .runId(runId)
                            .status(json.getString("status"))
                            .workflowId(workflowId)
                            .build();

                    System.out.println("RestWorkflowRunClient: Mapped RunResponse: id=" + runResponse.getRunId());
                    return runResponse;
                });
    }

    // Implement other methods similarly...

    @Override
    public Uni<RunResponse> getRun(String runId) {
        return webClient.get("/api/v1/workflow-runs/" + runId)
                .putHeader("X-Tenant-ID", config.tenantId())
                .putHeader("Authorization", "Bearer " + config.apiKey())
                .send()
                .map(response -> response.bodyAsJson(RunResponse.class));
    }

    // ... (other methods)

    @Override
    public Uni<RunResponse> startRun(String runId) {
        return webClient.post("/api/v1/workflow-runs/" + runId + "/start")
                .putHeader("X-Tenant-ID", config.tenantId())
                .putHeader("Authorization", "Bearer " + config.apiKey())
                .send()
                .map(response -> response.bodyAsJson(RunResponse.class));
    }

    @Override
    public Uni<RunResponse> suspendRun(String runId, String reason, String waitingOnNodeId) {
        return null;
    }

    @Override
    public Uni<RunResponse> resumeRun(String runId, Map<String, Object> resumeData, String humanTaskId) {
        return null;
    }

    @Override
    public Uni<Void> cancelRun(String runId, String reason) {
        return null;
    }

    @Override
    public Uni<Void> signal(String runId, String signalName, String targetNodeId, Map<String, Object> payload) {
        return null;
    }

    @Override
    public Uni<ExecutionHistory> getExecutionHistory(String runId) {
        return webClient.get("/api/v1/workflow-runs/" + runId + "/history")
                .putHeader("X-Tenant-ID", config.tenantId())
                .putHeader("Authorization", "Bearer " + config.apiKey())
                .send()
                .map(response -> response.bodyAsJson(ExecutionHistory.class));
    }

    @Override
    public Uni<List<RunResponse>> queryRuns(String workflowId, String status, int page, int size) {
        return null;
    }

    @Override
    public Uni<Long> getActiveRunsCount() {
        return null;
    }

    @Override
    public void close() {
        if (webClient != null) {
            webClient.close();
        }
        if (vertx != null) {
            vertx.close();
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/ResumeRunBuilder.java
Size: 1022 B | Modified: 2026-01-03 10:47:32
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.util.HashMap;
import java.util.Map;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.RunResponse;

/**
 * Builder for resuming runs
 */
public class ResumeRunBuilder {

    private final WorkflowRunClient client;
    private final String runId;
    private final Map<String, Object> resumeData = new HashMap<>();
    private String humanTaskId;

    ResumeRunBuilder(WorkflowRunClient client, String runId) {
        this.client = client;
        this.runId = runId;
    }

    public ResumeRunBuilder data(String key, Object value) {
        resumeData.put(key, value);
        return this;
    }

    public ResumeRunBuilder data(Map<String, Object> data) {
        this.resumeData.putAll(data);
        return this;
    }

    public ResumeRunBuilder humanTaskId(String taskId) {
        this.humanTaskId = taskId;
        return this;
    }

    public Uni<RunResponse> execute() {
        return client.resumeRun(runId, resumeData, humanTaskId);
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/SignalBuilder.java
Size: 1.1 KB | Modified: 2026-01-03 09:42:28
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.util.HashMap;
import java.util.Map;

import io.smallrye.mutiny.Uni;

/**
 * Builder for sending signals
 */
public class SignalBuilder {

    private final WorkflowRunClient client;
    private final String runId;
    private String signalName;
    private String targetNodeId;
    private final Map<String, Object> payload = new HashMap<>();

    SignalBuilder(WorkflowRunClient client, String runId) {
        this.client = client;
        this.runId = runId;
    }

    public SignalBuilder name(String signalName) {
        this.signalName = signalName;
        return this;
    }

    public SignalBuilder targetNode(String nodeId) {
        this.targetNodeId = nodeId;
        return this;
    }

    public SignalBuilder payload(String key, Object value) {
        payload.put(key, value);
        return this;
    }

    public SignalBuilder payload(Map<String, Object> payload) {
        this.payload.putAll(payload);
        return this;
    }

    public Uni<Void> send() {
        return client.signal(runId, signalName, targetNodeId, payload);
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/SilatClient.java
Size: 4.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.time.Duration;
import java.util.*;
import java.util.concurrent.atomic.AtomicBoolean;

/**
 * ============================================================================
 * SILAT CLIENT SDK
 * ============================================================================
 */
public class SilatClient implements AutoCloseable {

    private final SilatClientConfig config;
    private final io.vertx.mutiny.core.Vertx vertx;
    private final WorkflowRunClient runClient;
    private final WorkflowDefinitionClient definitionClient;
    private final AtomicBoolean closed = new AtomicBoolean(false);

    private SilatClient(SilatClientConfig config) {
        this.config = config;
        this.vertx = io.vertx.mutiny.core.Vertx.vertx();

        // Initialize transport-specific clients
        if (config.transport() == TransportType.REST) {
            this.runClient = new RestWorkflowRunClient(config, vertx);
            this.definitionClient = new RestWorkflowDefinitionClient(config, vertx);
        } else if (config.transport() == TransportType.GRPC) {
            this.runClient = new GrpcWorkflowRunClient(config);
            this.definitionClient = new GrpcWorkflowDefinitionClient(config);
        } else {
            throw new IllegalArgumentException("Unsupported transport: " + config.transport());
        }
    }

    /**
     * Get the client configuration
     */
    public SilatClientConfig config() {
        return config;
    }

    // ==================== BUILDER ====================

    public static Builder builder() {
        return new Builder();
    }

    public static class Builder {
        private String endpoint;
        private String tenantId;
        private String apiKey;
        private TransportType transport = TransportType.REST;
        private Duration timeout = Duration.ofSeconds(30);
        private Map<String, String> headers = new HashMap<>();

        public Builder restEndpoint(String endpoint) {
            this.endpoint = endpoint;
            this.transport = TransportType.REST;
            return this;
        }

        public Builder grpcEndpoint(String host, int port) {
            this.endpoint = host + ":" + port;
            this.transport = TransportType.GRPC;
            return this;
        }

        public Builder tenantId(String tenantId) {
            this.tenantId = tenantId;
            return this;
        }

        public Builder apiKey(String apiKey) {
            this.apiKey = apiKey;
            return this;
        }

        public Builder timeout(Duration timeout) {
            this.timeout = timeout;
            return this;
        }

        public Builder header(String key, String value) {
            this.headers.put(key, value);
            return this;
        }

        public SilatClient build() {
            SilatClientConfig config = SilatClientConfig.builder()
                    .endpoint(endpoint)
                    .tenantId(tenantId)
                    .apiKey(apiKey)
                    .transport(transport)
                    .timeout(timeout)
                    .headers(headers)
                    .build();

            return new SilatClient(config);
        }
    }

    // ==================== API METHODS ====================

    /**
     * Access workflow run operations
     */
    public WorkflowRunOperations runs() {
        checkClosed();
        return new WorkflowRunOperations(runClient);
    }

    /**
     * Access workflow definition operations
     */
    public WorkflowDefinitionOperations workflows() {
        checkClosed();
        return new WorkflowDefinitionOperations(definitionClient);
    }

    private void checkClosed() {
        if (closed.get()) {
            throw new IllegalStateException("SilatClient is closed");
        }
    }

    /**
     * Close the client and release resources
     */
    @Override
    public void close() {
        if (closed.compareAndSet(false, true)) {
            if (runClient != null) {
                runClient.close();
            }
            if (definitionClient != null) {
                definitionClient.close();
            }
            if (vertx != null) {
                vertx.close();
            }
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/SilatClientConfig.java
Size: 4.4 KB | Modified: 2026-01-03 17:06:55
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.time.Duration;
import java.util.Collections;
import java.util.Map;
import java.util.Objects;

/**
 * Configuration for the Silat client.
 * This class holds all the necessary configuration parameters for connecting to
 * the Silat service.
 */
public final class SilatClientConfig {
    private final String endpoint;
    private final String tenantId;
    private final String apiKey;
    private final TransportType transport;
    private final Duration timeout;
    private final Map<String, String> headers;

    private SilatClientConfig(String endpoint, String tenantId, String apiKey,
            TransportType transport, Duration timeout, Map<String, String> headers) {
        this.endpoint = endpoint;
        this.tenantId = tenantId;
        this.apiKey = apiKey;
        this.transport = transport;
        this.timeout = timeout;
        this.headers = headers != null ? Collections.unmodifiableMap(headers) : Map.of();
    }

    // Getters
    public String endpoint() {
        return endpoint;
    }

    public String tenantId() {
        return tenantId;
    }

    public String apiKey() {
        return apiKey;
    }

    public TransportType transport() {
        return transport;
    }

    public Duration timeout() {
        return timeout;
    }

    public Map<String, String> headers() {
        return headers;
    }

    public static Builder builder() {
        return new Builder();
    }

    public static SilatClientConfig defaultConfig(String endpoint, String tenantId) {
        return builder()
                .endpoint(endpoint)
                .tenantId(tenantId)
                .transport(TransportType.REST)
                .timeout(Duration.ofSeconds(30))
                .build();
    }

    public static class Builder {
        private String endpoint;
        private String tenantId;
        private String apiKey;
        private TransportType transport = TransportType.REST;
        private Duration timeout = Duration.ofSeconds(30);
        private Map<String, String> headers = new java.util.HashMap<>();

        public Builder endpoint(String endpoint) {
            this.endpoint = endpoint;
            return this;
        }

        public Builder tenantId(String tenantId) {
            this.tenantId = tenantId;
            return this;
        }

        public Builder apiKey(String apiKey) {
            this.apiKey = apiKey;
            return this;
        }

        public Builder transport(TransportType transport) {
            this.transport = transport;
            return this;
        }

        public Builder timeout(Duration timeout) {
            this.timeout = timeout;
            return this;
        }

        public Builder header(String key, String value) {
            this.headers.put(key, value);
            return this;
        }

        public Builder headers(Map<String, String> headers) {
            this.headers.putAll(headers);
            return this;
        }

        public SilatClientConfig build() {
            Objects.requireNonNull(endpoint, "Endpoint cannot be null");
            Objects.requireNonNull(tenantId, "Tenant ID cannot be null");
            Objects.requireNonNull(transport, "Transport type cannot be null");
            Objects.requireNonNull(timeout, "Timeout cannot be null");

            if (endpoint.trim().isEmpty()) {
                throw new IllegalArgumentException("Endpoint cannot be empty");
            }
            if (tenantId.trim().isEmpty()) {
                throw new IllegalArgumentException("Tenant ID cannot be empty");
            }
            if (timeout.isNegative() || timeout.isZero()) {
                throw new IllegalArgumentException("Timeout must be positive");
            }
            if (apiKey != null && apiKey.trim().isEmpty()) {
                throw new IllegalArgumentException("API key cannot be empty when provided");
            }

            return new SilatClientConfig(endpoint, tenantId, apiKey, transport, timeout, headers);
        }

        public Builder rest() {
            this.transport = TransportType.REST;
            return this;
        }

        public Builder grpc() {
            this.transport = TransportType.GRPC;
            return this;
        }

        public Builder timeoutSeconds(long seconds) {
            this.timeout = Duration.ofSeconds(seconds);
            return this;
        }
    }
}
================================================================================

================================================================================
tech/kayys/silat/sdk/client/SuspendRunBuilder.java
Size: 807 B | Modified: 2026-01-03 10:40:14
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.RunResponse;

/**
 * Builder for suspending runs
 */
public class SuspendRunBuilder {

    private final WorkflowRunClient client;
    private final String runId;
    private String reason;
    private String waitingOnNodeId;

    SuspendRunBuilder(WorkflowRunClient client, String runId) {
        this.client = client;
        this.runId = runId;
    }

    public SuspendRunBuilder reason(String reason) {
        this.reason = reason;
        return this;
    }

    public SuspendRunBuilder waitingOnNode(String nodeId) {
        this.waitingOnNodeId = nodeId;
        return this;
    }

    public Uni<RunResponse> execute() {
        return client.suspendRun(runId, reason, waitingOnNodeId);
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/TransportType.java
Size: 238 B | Modified: 2026-01-03 10:38:43
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

/**
 * Transport protocol type for the Silat client.
 */
public enum TransportType {
    /**
     * REST transport protocol
     */
    REST,

    /**
     * gRPC transport protocol
     */
    GRPC
}
================================================================================

================================================================================
tech/kayys/silat/sdk/client/WorkflowDefinitionBuilder.java
Size: 3.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.time.Instant;
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.WorkflowDefinition;
import tech.kayys.silat.model.WorkflowDefinitionId;
import tech.kayys.silat.model.TenantId;
import tech.kayys.silat.model.WorkflowMetadata;
import tech.kayys.silat.model.NodeDefinition;
import tech.kayys.silat.model.InputDefinition;
import tech.kayys.silat.model.OutputDefinition;
import tech.kayys.silat.model.RetryPolicy;
import tech.kayys.silat.saga.CompensationPolicy;

/**
 * Builder for workflow definitions
 */
public class WorkflowDefinitionBuilder {

    private final WorkflowDefinitionClient client;
    private final String name;
    private String version = "1.0.0";
    private String tenantId = "default";
    private String description;
    private final List<NodeDefinition> nodes = new ArrayList<>();
    private final Map<String, InputDefinition> inputs = new HashMap<>();
    private final Map<String, OutputDefinition> outputs = new HashMap<>();
    private RetryPolicy retryPolicy;
    private CompensationPolicy compensationPolicy;
    private final Map<String, String> labels = new HashMap<>();

    WorkflowDefinitionBuilder(WorkflowDefinitionClient client, String name) {
        this.client = client;
        this.name = name;
    }

    public WorkflowDefinitionBuilder version(String version) {
        this.version = version;
        return this;
    }

    public WorkflowDefinitionBuilder tenantId(String tenantId) {
        this.tenantId = tenantId;
        return this;
    }

    public WorkflowDefinitionBuilder description(String description) {
        this.description = description;
        return this;
    }

    public WorkflowDefinitionBuilder addNode(NodeDefinition node) {
        nodes.add(node);
        return this;
    }

    public WorkflowDefinitionBuilder addInput(String name, InputDefinition input) {
        inputs.put(name, input);
        return this;
    }

    public WorkflowDefinitionBuilder addOutput(String name, OutputDefinition output) {
        outputs.put(name, output);
        return this;
    }

    public WorkflowDefinitionBuilder retryPolicy(RetryPolicy policy) {
        this.retryPolicy = policy;
        return this;
    }

    public WorkflowDefinitionBuilder compensationPolicy(CompensationPolicy policy) {
        this.compensationPolicy = policy;
        return this;
    }

    public WorkflowDefinitionBuilder label(String key, String value) {
        labels.put(key, value);
        return this;
    }

    public Uni<WorkflowDefinition> execute() {
        WorkflowMetadata metadata = new WorkflowMetadata(
                labels,
                new HashMap<>(), // annotations
                Instant.now(),
                "sdk-client");

        WorkflowDefinition request = new WorkflowDefinition(
                WorkflowDefinitionId.of(name),
                TenantId.of(tenantId),
                name,
                version,
                description,
                nodes,
                inputs,
                outputs,
                metadata,
                retryPolicy,
                compensationPolicy);
        return client.createDefinition(request);
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/WorkflowDefinitionClient.java
Size: 550 B | Modified: 2026-01-03 17:37:50
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.util.List;
import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.WorkflowDefinition;

/**
 * Workflow definition client interface
 */
interface WorkflowDefinitionClient extends AutoCloseable {
    Uni<WorkflowDefinition> createDefinition(WorkflowDefinition request);

    Uni<WorkflowDefinition> getDefinition(String definitionId);

    Uni<List<WorkflowDefinition>> listDefinitions(boolean activeOnly);

    Uni<Void> deleteDefinition(String definitionId);

    @Override
    void close();
}

================================================================================

================================================================================
tech/kayys/silat/sdk/client/WorkflowDefinitionOperations.java
Size: 1.0 KB | Modified: 2026-01-03 10:41:04
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.util.List;
import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.WorkflowDefinition;

/**
 * Fluent API for workflow definition operations
 */
public class WorkflowDefinitionOperations {

    private final WorkflowDefinitionClient client;

    WorkflowDefinitionOperations(WorkflowDefinitionClient client) {
        this.client = client;
    }

    /**
     * Create a new workflow definition
     */
    public WorkflowDefinitionBuilder create(String name) {
        return new WorkflowDefinitionBuilder(client, name);
    }

    /**
     * Get a workflow definition
     */
    public Uni<WorkflowDefinition> get(String definitionId) {
        return client.getDefinition(definitionId);
    }

    /**
     * List workflow definitions
     */
    public Uni<List<WorkflowDefinition>> list() {
        return client.listDefinitions(true);
    }

    /**
     * Delete a workflow definition
     */
    public Uni<Void> delete(String definitionId) {
        return client.deleteDefinition(definitionId);
    }
}
================================================================================

================================================================================
tech/kayys/silat/sdk/client/WorkflowRunClient.java
Size: 1.1 KB | Modified: 2026-01-03 17:38:27
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import java.util.Map;
import java.util.List;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.CreateRunRequest;
import tech.kayys.silat.model.RunResponse;
import tech.kayys.silat.execution.ExecutionHistory;

/**
 * Workflow run client interface (transport-agnostic)
 */
interface WorkflowRunClient extends AutoCloseable {
    Uni<RunResponse> createRun(CreateRunRequest request);

    Uni<RunResponse> getRun(String runId);

    Uni<RunResponse> startRun(String runId);

    Uni<RunResponse> suspendRun(String runId, String reason, String waitingOnNodeId);

    Uni<RunResponse> resumeRun(String runId, Map<String, Object> resumeData, String humanTaskId);

    Uni<Void> cancelRun(String runId, String reason);

    Uni<Void> signal(String runId, String signalName, String targetNodeId, Map<String, Object> payload);

    Uni<ExecutionHistory> getExecutionHistory(String runId);

    Uni<List<RunResponse>> queryRuns(String workflowId, String status, int page, int size);

    Uni<Long> getActiveRunsCount();

    @Override
    void close();
}
================================================================================

================================================================================
tech/kayys/silat/sdk/client/WorkflowRunOperations.java
Size: 1.9 KB | Modified: 2026-01-03 17:40:33
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.client;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.model.RunResponse;
import tech.kayys.silat.execution.ExecutionHistory;

/**
 * Fluent API for workflow run operations
 */
public class WorkflowRunOperations {

    private final WorkflowRunClient client;

    WorkflowRunOperations(WorkflowRunClient client) {
        this.client = client;
    }

    /**
     * Create a new workflow run
     */
    public CreateRunBuilder create(String workflowDefinitionId) {
        return new CreateRunBuilder(client, workflowDefinitionId);
    }

    /**
     * Get a workflow run
     */
    public Uni<RunResponse> get(String runId) {
        return client.getRun(runId);
    }

    /**
     * Start a workflow run
     */
    public Uni<RunResponse> start(String runId) {
        return client.startRun(runId);
    }

    /**
     * Suspend a workflow run
     */
    public SuspendRunBuilder suspend(String runId) {
        return new SuspendRunBuilder(client, runId);
    }

    /**
     * Resume a workflow run
     */
    public ResumeRunBuilder resume(String runId) {
        return new ResumeRunBuilder(client, runId);
    }

    /**
     * Cancel a workflow run
     */
    public Uni<Void> cancel(String runId, String reason) {
        return client.cancelRun(runId, reason);
    }

    /**
     * Send signal to workflow run
     */
    public SignalBuilder signal(String runId) {
        return new SignalBuilder(client, runId);
    }

    /**
     * Get execution history
     */
    public Uni<ExecutionHistory> getHistory(String runId) {
        return client.getExecutionHistory(runId);
    }

    /**
     * Query workflow runs
     */
    public QueryRunsBuilder query() {
        return new QueryRunsBuilder(client);
    }

    /**
     * Get active runs count
     */
    public Uni<Long> getActiveCount() {
        return client.getActiveRunsCount();
    }
}

================================================================================

Coto Output
Generated: 2026-01-17 09:25:48
Files: 16 | Directories: 6 | Total Size: 54.1 KB


================================================================================
tech/kayys/silat/sdk/executor/AbstractWorkflowExecutor.java
Size: 8.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.time.Duration;
import java.time.Instant;
import java.util.Arrays;
import java.util.concurrent.atomic.AtomicInteger;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;
import tech.kayys.silat.model.ErrorInfo;

/**
 * Abstract base class for executors with common functionality
 */
public abstract class AbstractWorkflowExecutor implements WorkflowExecutor {

    private static final Logger LOG = LoggerFactory.getLogger(AbstractWorkflowExecutor.class);

    protected final String executorType;
    protected final ExecutorConfig config;
    protected final ExecutorMetrics metrics;
    protected final AtomicInteger activeTaskCount = new AtomicInteger(0);

    protected AbstractWorkflowExecutor() {
        // Extract executor type from annotation
        Executor annotation = getClass().getAnnotation(Executor.class);

        // Handle Quarkus/CDI proxies
        if (annotation == null && getClass().getName().endsWith("_Subclass")) {
            annotation = getClass().getSuperclass().getAnnotation(Executor.class);
        }

        if (annotation == null) {
            LOG.error("Failed to find @Executor annotation on class: {}", getClass().getName());
            // Attempt to look up the hierarchy
            Class<?> current = getClass();
            while (current != Object.class) {
                annotation = current.getAnnotation(Executor.class);
                if (annotation != null) {
                    LOG.info("Found @Executor annotation on parent class: {}", current.getName());
                    break;
                }
                current = current.getSuperclass();
            }
        }

        if (annotation == null) {
            throw new IllegalStateException(
                    "Executor class must be annotated with @Executor: " + getClass().getName());
        }

        this.executorType = annotation.executorType();
        this.config = new ExecutorConfig(
                annotation.maxConcurrentTasks(),
                Arrays.asList(annotation.supportedNodeTypes()),
                annotation.communicationType(),
                SecurityConfig.disabled());
        this.metrics = new ExecutorMetrics(executorType);
    }

    @Override
    public final String getExecutorType() {
        return executorType;
    }

    @Override
    public int getMaxConcurrentTasks() {
        return config.maxConcurrentTasks();
    }

    @Override
    public String[] getSupportedNodeTypes() {
        return config.supportedNodeTypes().toArray(new String[0]);
    }

    @Override
    public boolean isReady() {
        return activeTaskCount.get() < getMaxConcurrentTasks();
    }

    /**
     * Execute with comprehensive lifecycle hooks and error handling
     */
    public final Uni<NodeExecutionResult> executeWithLifecycle(NodeExecutionTask task) {
        LOG.debug("Executing task: run={}, node={}, attempt={}, executor={}",
                task.runId().value(), task.nodeId().value(), task.attempt(), executorType);

        // Check if executor is ready to handle the task
        if (!isReady()) {
            LOG.warn("Executor {} is not ready, active tasks: {}, max: {}",
                    executorType, activeTaskCount.get(), getMaxConcurrentTasks());
            return Uni.createFrom().item(SimpleNodeExecutionResult.failure(
                    task.runId(),
                    task.nodeId(),
                    task.attempt(),
                    ErrorInfo.of(new IllegalStateException("Executor not ready - too many active tasks")),
                    task.token()));
        }

        Instant startTime = Instant.now();
        activeTaskCount.incrementAndGet();
        metrics.recordTaskStarted();

        return beforeExecute(task)
                .onItem().invoke(() -> LOG.trace("Before execute completed for task: {}", task.nodeId()))
                .onFailure()
                .invoke(throwable -> LOG.error("Before execute failed for task: {}", task.nodeId(), throwable))
                .flatMap(v -> execute(task))
                .onItem().invoke(result -> {
                    Duration duration = Duration.between(startTime, Instant.now());
                    activeTaskCount.decrementAndGet();
                    metrics.recordTaskCompleted(duration);
                    LOG.info("Task completed: run={}, node={}, status={}, duration={}ms, attempt={}",
                            task.runId().value(), task.nodeId().value(),
                            result.status(), duration.toMillis(), task.attempt());
                })
                .onFailure().invoke(throwable -> {
                    Duration duration = Duration.between(startTime, Instant.now());
                    activeTaskCount.decrementAndGet();
                    metrics.recordTaskFailed(duration);
                    LOG.error("Task failed: run={}, node={}, attempt={}",
                            task.runId().value(), task.nodeId().value(), task.attempt(), throwable);

                    // Call onError hook for error handling
                    onError(task, throwable)
                            .subscribe().with(
                                    ignored -> LOG.trace("onError hook completed for task: {}", task.nodeId()),
                                    error -> LOG.warn("onError hook failed for task: {}", task.nodeId(), error));
                })
                .onFailure().recoverWithItem(throwable -> {
                    Duration duration = Duration.between(startTime, Instant.now());
                    activeTaskCount.decrementAndGet();
                    metrics.recordTaskFailed(duration);
                    LOG.warn("Recovering from execution failure for task: {}", task.nodeId(), throwable);
                    return SimpleNodeExecutionResult.failure(
                            task.runId(),
                            task.nodeId(),
                            task.attempt(),
                            ErrorInfo.of(throwable),
                            task.token());
                })
                .flatMap(result -> afterExecute(task, result)
                        .onItem().invoke(() -> LOG.trace("After execute completed for task: {}", task.nodeId()))
                        .onFailure()
                        .invoke(throwable -> LOG.warn("After execute failed for task: {}", task.nodeId(), throwable))
                        .replaceWith(result));
    }

    /**
     * Validates if the executor can handle the given task
     */
    @Override
    public boolean canHandle(NodeExecutionTask task) {
        // Check if the node type is supported
        String[] supportedTypes = getSupportedNodeTypes();
        if (supportedTypes.length > 0) {
            String nodeType = extractNodeType(task);
            for (String supportedType : supportedTypes) {
                if (supportedType.equals(nodeType)) {
                    return true;
                }
            }
            LOG.debug("Executor {} cannot handle node type: {} for task: {}",
                    executorType, nodeType, task.nodeId());
            return false;
        }
        return true; // If no specific types defined, assume it can handle any
    }

    /**
     * Extracts node type from the task (implementation may vary based on actual
     * task structure)
     */
    protected String extractNodeType(NodeExecutionTask task) {
        // Look for the special __node_type__ key in the context, which is the system
        // convention
        if (task.context() != null && task.context().containsKey("__node_type__")) {
            return String.valueOf(task.context().get("__node_type__"));
        }

        // Fallback to node ID value if not found
        return task.nodeId().value();
    }

    /**
     * Gets the current number of active tasks
     */
    public int getActiveTaskCount() {
        return activeTaskCount.get();
    }

    /**
     * Gets the configuration for this executor
     */
    public ExecutorConfig getConfig() {
        return config;
    }

    /**
     * Gets the metrics for this executor
     */
    public ExecutorMetrics getMetrics() {
        return metrics;
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/Executor.java
Size: 1.9 KB | Modified: 2026-01-03 18:02:39
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import tech.kayys.silat.model.CommunicationType;

import java.lang.annotation.*;

/**
 * ============================================================================
 * SILAT EXECUTOR SDK
 * ============================================================================
 * 
 * Framework for building workflow executors that process node tasks.
 * Supports multiple communication strategies (gRPC, Kafka, REST).
 * 
 * Example Usage:
 * ```java
 * @Executor(
 *     executorType = "order-validator",
 *     communicationType = CommunicationType.GRPC
 * )
 * public class OrderValidatorExecutor extends AbstractWorkflowExecutor {
 *     
 *     @Override
 *     public Uni<NodeExecutionResult> execute(NodeExecutionTask task) {
 *         Map<String, Object> context = task.context();
 *         String orderId = (String) context.get("orderId");
 *         
 *         return validateOrder(orderId)
 *             .map(valid -> NodeExecutionResult.success(
 *                 task.runId(),
 *                 task.nodeId(),
 *                 task.attempt(),
 *                 Map.of("valid", valid),
 *                 task.token()
 *             ));
 *     }
 * }
 * ```
 */

// ==================== EXECUTOR ANNOTATION ====================

/**
 * Annotation to mark a class as a workflow executor
 */
@Target(ElementType.TYPE)
@Retention(RetentionPolicy.RUNTIME)
@Documented
public @interface Executor {

    /**
     * Unique executor type identifier
     */
    String executorType();

    /**
     * Communication type for receiving tasks
     */
    CommunicationType communicationType() default CommunicationType.GRPC;

    /**
     * Maximum concurrent tasks
     */
    int maxConcurrentTasks() default 10;

    /**
     * Supported node types
     */
    String[] supportedNodeTypes() default {};

    /**
     * Executor version
     */
    String version() default "1.0.0";
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/ExecutorConfig.java
Size: 363 B | Modified: 2026-01-06 13:27:14
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.util.List;

import tech.kayys.silat.model.CommunicationType;

/**
 * Executor configuration
 */
record ExecutorConfig(
                int maxConcurrentTasks,
                List<String> supportedNodeTypes,
                CommunicationType communicationType,
                SecurityConfig securityConfig) {
}
================================================================================

================================================================================
tech/kayys/silat/sdk/executor/ExecutorMetrics.java
Size: 1.6 KB | Modified: 2026-01-03 18:00:26
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.time.Duration;
import java.util.List;
import java.util.Map;
import java.util.concurrent.CopyOnWriteArrayList;

/**
 * Executor metrics
 */
class ExecutorMetrics {

    private final String executorType;
    private final java.util.concurrent.atomic.AtomicLong tasksStarted = new java.util.concurrent.atomic.AtomicLong();
    private final java.util.concurrent.atomic.AtomicLong tasksCompleted = new java.util.concurrent.atomic.AtomicLong();
    private final java.util.concurrent.atomic.AtomicLong tasksFailed = new java.util.concurrent.atomic.AtomicLong();
    private final List<Duration> durations = new CopyOnWriteArrayList<>();

    ExecutorMetrics(String executorType) {
        this.executorType = executorType;
    }

    void recordTaskStarted() {
        tasksStarted.incrementAndGet();
    }

    void recordTaskCompleted(Duration duration) {
        tasksCompleted.incrementAndGet();
        durations.add(duration);
    }

    void recordTaskFailed(Duration duration) {
        tasksFailed.incrementAndGet();
        durations.add(duration);
    }

    public Map<String, Object> getMetrics() {
        return Map.of(
                "executorType", executorType,
                "tasksStarted", tasksStarted.get(),
                "tasksCompleted", tasksCompleted.get(),
                "tasksFailed", tasksFailed.get(),
                "avgDurationMs", calculateAvgDuration());
    }

    private long calculateAvgDuration() {
        if (durations.isEmpty())
            return 0;
        return durations.stream()
                .mapToLong(Duration::toMillis)
                .sum() / durations.size();
    }
}
================================================================================

================================================================================
tech/kayys/silat/sdk/executor/ExecutorPluginManager.java
Size: 2.3 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.util.Comparator;
import java.util.List;
import java.util.stream.Collectors;

import jakarta.annotation.PostConstruct;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import tech.kayys.silat.execution.NodeExecutionTask;
import tech.kayys.silat.plugin.executor.ExecutorPlugin;
import tech.kayys.silat.plugin.impl.PluginManager;

/**
 * Manager for executor plugins
 * 
 * Discovers and manages executor plugins, providing plugin selection
 * based on task requirements.
 */
@ApplicationScoped
public class ExecutorPluginManager {

    private static final Logger LOG = LoggerFactory.getLogger(ExecutorPluginManager.class);

    @Inject
    PluginManager pluginManager;

    private List<ExecutorPlugin> executorPlugins;

    @PostConstruct
    void init() {
        // Discover and load executor plugins
        executorPlugins = pluginManager.getPluginsByType(ExecutorPlugin.class).stream()
                .sorted(Comparator.comparingInt(ExecutorPlugin::getPriority).reversed())
                .collect(Collectors.toList());

        LOG.info("Loaded {} executor plugins", executorPlugins.size());
        executorPlugins.forEach(p -> LOG.info("  - {} (type: {}, priority: {})",
                p.getMetadata().name(),
                p.getExecutorType(),
                p.getPriority()));
    }

    /**
     * Find a suitable plugin for the given task
     * 
     * @param task the task to find a plugin for
     * @return the first plugin that can handle the task, or null if none found
     */
    public ExecutorPlugin findPlugin(NodeExecutionTask task) {
        return executorPlugins.stream()
                .filter(p -> p.canHandle(task))
                .findFirst()
                .orElse(null);
    }

    /**
     * Check if any plugin can handle the given task
     * 
     * @param task the task to check
     * @return true if at least one plugin can handle the task
     */
    public boolean hasPluginFor(NodeExecutionTask task) {
        return findPlugin(task) != null;
    }

    /**
     * Get all loaded executor plugins
     * 
     * @return list of executor plugins
     */
    public List<ExecutorPlugin> getPlugins() {
        return executorPlugins;
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/ExecutorRuntime.java
Size: 6.0 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.util.ArrayList;
import java.util.Map;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.quarkus.runtime.Startup;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import jakarta.enterprise.context.ApplicationScoped;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;
import tech.kayys.silat.model.ErrorInfo;

/**
 * Runtime for managing executor lifecycle
 */
@Startup
@ApplicationScoped
public class ExecutorRuntime {

    private static final Logger LOG = LoggerFactory.getLogger(ExecutorRuntime.class);

    private final Map<String, WorkflowExecutor> executors = new ConcurrentHashMap<>();
    private final ExecutorService executorService;
    private final ExecutorTransport transport;
    private volatile boolean running = false;

    @jakarta.inject.Inject
    public ExecutorRuntime(ExecutorTransportFactory transportFactory,
            jakarta.enterprise.inject.Instance<WorkflowExecutor> discoveredExecutors) {
        this.executorService = Executors.newVirtualThreadPerTaskExecutor();
        this.transport = transportFactory.createTransport();

        // Auto-discover and register executors
        discoveredExecutors.forEach(executor -> {
            String type = executor.getExecutorType();
            executors.put(type, executor);
            LOG.info("Auto-discovered and registered executor: {}", type);
        });
    }

    /**
     * Register an executor manually if needed
     */
    public void registerExecutor(WorkflowExecutor executor) {
        String type = executor.getExecutorType();
        executors.put(type, executor);
        LOG.info("Manually registered executor: {}", type);
    }

    /**
     * Start the runtime
     */
    @PostConstruct
    public void start() {
        LOG.info("Starting executor runtime with {} executors", executors.size());
        running = true;

        // Register with engine
        // Register with engine - delay to ensure listener is ready
        io.smallrye.mutiny.Uni.createFrom().voidItem()
                .onItem().delayIt().by(java.time.Duration.ofSeconds(2))
                .flatMap(v -> transport.register(new ArrayList<>(executors.values())))
                .subscribe().with(
                        v -> LOG.info("Registered with engine"),
                        error -> LOG.error("Failed to register", error));

        // Start receiving tasks
        transport.receiveTasks()
                .subscribe().with(
                        task -> handleTask(task),
                        error -> LOG.error("Error receiving tasks", error));

        // Start heartbeat
        startHeartbeat();
    }

    /**
     * Stop the runtime
     */
    @PreDestroy
    public void stop() {
        LOG.info("Stopping executor runtime");
        running = false;

        // Unregister from engine
        transport.unregister()
                .subscribe().with(
                        v -> LOG.info("Unregistered from engine"),
                        error -> LOG.error("Failed to unregister", error));

        executorService.shutdown();
    }

    /**
     * Handle incoming task
     */
    private void handleTask(NodeExecutionTask task) {
        LOG.debug("Received task: run={}, node={}",
                task.runId().value(), task.nodeId().value());

        // Find appropriate executor
        WorkflowExecutor executor = executors.values().stream()
                .filter(e -> e.canHandle(task))
                .findFirst()
                .orElse(null);

        if (executor == null) {
            LOG.warn("No executor found for task: {}", task.nodeId().value());
            sendResult(SimpleNodeExecutionResult.failure(
                    task.runId(),
                    task.nodeId(),
                    task.attempt(),
                    new ErrorInfo("NO_EXECUTOR", "No executor found", "", Map.of()),
                    task.token()));
            return;
        }

        // Execute in virtual thread
        executorService.submit(() -> {
            if (executor instanceof AbstractWorkflowExecutor abstractExecutor) {
                abstractExecutor.executeWithLifecycle(task)
                        .subscribe().with(
                                result -> sendResult(result),
                                error -> LOG.error("Execution failed", error));
            } else {
                executor.execute(task)
                        .subscribe().with(
                                result -> sendResult(result),
                                error -> LOG.error("Execution failed", error));
            }
        });
    }

    /**
     * Send result back to engine
     */
    private void sendResult(NodeExecutionResult result) {
        LOG.debug("Sending result: run={}, node={}, status={}",
                result.runId().value(), result.nodeId().value(), result.status());

        transport.sendResult(result)
                .subscribe().with(
                        v -> LOG.debug("Result sent successfully"),
                        error -> LOG.error("Failed to send result", error));
    }

    /**
     * Send periodic heartbeat
     */
    private void startHeartbeat() {
        CompletableFuture.runAsync(() -> {
            while (running) {
                try {
                    transport.sendHeartbeat()
                            .subscribe().with(
                                    v -> LOG.trace("Heartbeat sent"),
                                    error -> LOG.warn("Heartbeat failed", error));

                    Thread.sleep(transport.getHeartbeatInterval().toMillis());
                } catch (InterruptedException e) {
                    Thread.currentThread().interrupt();
                    break;
                }
            }
        }, executorService);
    }
}
================================================================================

================================================================================
tech/kayys/silat/sdk/executor/ExecutorTransport.java
Size: 1.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.util.List;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;

/**
 * Transport interface for executor communication
 */
public interface ExecutorTransport {

    /**
     * Register executors with engine
     */
    Uni<Void> register(List<WorkflowExecutor> executors);

    /**
     * Unregister from engine
     */
    Uni<Void> unregister();

    /**
     * Receive tasks from engine (streaming)
     */
    io.smallrye.mutiny.Multi<NodeExecutionTask> receiveTasks();

    /**
     * Send task result to engine
     */
    Uni<Void> sendResult(NodeExecutionResult result);

    /**
     * Send heartbeat
     */
    Uni<Void> sendHeartbeat();

    /**
     * Get the communication type of this transport
     */
    default tech.kayys.silat.model.CommunicationType getCommunicationType() {
        return tech.kayys.silat.model.CommunicationType.UNSPECIFIED;
    }

    /**
     * Get configured heartbeat interval
     * 
     * @return Duration interval
     */
    default java.time.Duration getHeartbeatInterval() {
        return java.time.Duration.ofSeconds(30);
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/ExecutorTransportFactory.java
Size: 942 B | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

/**
 * Factory for creating executor transports
 */
@jakarta.enterprise.context.ApplicationScoped
public class ExecutorTransportFactory {

    @jakarta.inject.Inject
    jakarta.enterprise.inject.Instance<ExecutorTransport> availableTransports;

    @org.eclipse.microprofile.config.inject.ConfigProperty(name = "silat.executor.transport", defaultValue = "GRPC")
    String transportType;

    public ExecutorTransport createTransport() {
        for (ExecutorTransport transport : availableTransports) {
            if (transport.getCommunicationType().name().equalsIgnoreCase(transportType)) {
                return transport;
            }
        }

        throw new IllegalArgumentException(
                "Unknown or unavailable transport: " + transportType + ". Available: " +
                        availableTransports.stream().map(t -> t.getCommunicationType().name()).toList());
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/GrpcExecutorTransport.java
Size: 11.7 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.UUID;
import java.util.concurrent.Executor;
import java.util.concurrent.Executors;
import java.util.concurrent.ScheduledExecutorService;
import java.util.concurrent.TimeUnit;
import java.util.concurrent.atomic.AtomicBoolean;

import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.grpc.ConnectivityState;
import io.grpc.ManagedChannel;
import io.grpc.netty.GrpcSslContexts;
import io.grpc.netty.NettyChannelBuilder;
import io.netty.handler.ssl.SslContext;
import io.netty.handler.ssl.SslContextBuilder;
import io.smallrye.mutiny.Uni;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.operators.multi.processors.BroadcastProcessor;
import jakarta.annotation.PostConstruct;
import jakarta.annotation.PreDestroy;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;
import tech.kayys.silat.grpc.GrpcMapper;
import tech.kayys.silat.grpc.v1.MutinyExecutorServiceGrpc;
import tech.kayys.silat.grpc.v1.RegisterExecutorRequest;
import tech.kayys.silat.grpc.v1.UnregisterExecutorRequest;
import tech.kayys.silat.grpc.v1.HeartbeatRequest;
import tech.kayys.silat.grpc.v1.StreamTasksRequest;
import tech.kayys.silat.grpc.v1.TaskResult;
import tech.kayys.silat.model.ExecutionToken;
import tech.kayys.silat.model.NodeId;
import tech.kayys.silat.model.WorkflowRunId;

/**
 * gRPC-based executor transport
 */
@ApplicationScoped
public class GrpcExecutorTransport implements ExecutorTransport {

    private static final Logger LOG = LoggerFactory.getLogger(GrpcExecutorTransport.class);

    private final String executorId;

    @ConfigProperty(name = "engine.grpc.endpoint", defaultValue = "localhost")
    String engineEndpoint;

    @ConfigProperty(name = "engine.grpc.port", defaultValue = "9090")
    int grpcPort;

    @ConfigProperty(name = "heartbeat.interval", defaultValue = "30s")
    Duration heartbeatInterval;

    @ConfigProperty(name = "grpc.max.retries", defaultValue = "3")
    int maxRetries;

    @ConfigProperty(name = "grpc.retry.delay", defaultValue = "5s")
    Duration retryDelay;

    @ConfigProperty(name = "security.mtls.enabled", defaultValue = "false")
    boolean mtlsEnabled;

    @ConfigProperty(name = "security.jwt.enabled", defaultValue = "false")
    boolean jwtEnabled;

    @ConfigProperty(name = "security.mtls.cert.path")
    java.util.Optional<String> keyCertChainPath;

    @ConfigProperty(name = "security.mtls.key.path")
    java.util.Optional<String> privateKeyPath;

    @ConfigProperty(name = "security.mtls.trust.path")
    java.util.Optional<String> trustCertCollectionPath;

    @ConfigProperty(name = "security.jwt.token")
    java.util.Optional<String> jwtToken;

    @Inject
    GrpcMapper mapper;

    private ManagedChannel channel;
    private MutinyExecutorServiceGrpc.MutinyExecutorServiceStub stub;
    private final AtomicBoolean isConnected = new AtomicBoolean(false);
    private final AtomicBoolean isShutdown = new AtomicBoolean(false);

    // For streaming task reception
    private final BroadcastProcessor<NodeExecutionTask> taskProcessor = BroadcastProcessor.create();

    // For background operations
    private final Executor executor = Executors.newSingleThreadExecutor();
    private final ScheduledExecutorService scheduledExecutor = Executors.newScheduledThreadPool(2);

    public GrpcExecutorTransport() {
        this.executorId = UUID.randomUUID().toString();
    }

    @PostConstruct
    public void init() {
        initializeChannel();
    }

    private void initializeChannel() {
        NettyChannelBuilder channelBuilder = NettyChannelBuilder
                .forAddress(engineEndpoint, grpcPort)
                .keepAliveTime(1, TimeUnit.MINUTES)
                .keepAliveTimeout(20, TimeUnit.SECONDS)
                .keepAliveWithoutCalls(true)
                .maxInboundMessageSize(4 * 1024 * 1024) // 4MB
                .defaultLoadBalancingPolicy("round_robin");

        if (mtlsEnabled) {
            LOG.info("Configuring mTLS for gRPC channel");
            try {
                SslContextBuilder sslContextBuilder = GrpcSslContexts.forClient();
                if (trustCertCollectionPath.isPresent()) {
                    sslContextBuilder.trustManager(new java.io.File(trustCertCollectionPath.get()));
                }
                if (keyCertChainPath.isPresent() && privateKeyPath.isPresent()) {
                    sslContextBuilder.keyManager(
                            new java.io.File(keyCertChainPath.get()),
                            new java.io.File(privateKeyPath.get()));
                }
                SslContext sslContext = sslContextBuilder.build();
                channelBuilder.sslContext(sslContext).useTransportSecurity();
            } catch (Exception e) {
                LOG.error("Failed to configure mTLS", e);
            }
        } else {
            channelBuilder.usePlaintext();
        }

        if (jwtEnabled && jwtToken.isPresent()) {
            LOG.info("Configuring JWT interceptor for gRPC channel");
            // Note: JwtClientInterceptor needs to be available in classpath if used
            // channelBuilder.intercept(new JwtClientInterceptor(jwtToken.get()));
        }

        this.channel = channelBuilder.build();
        this.stub = MutinyExecutorServiceGrpc.newMutinyStub(channel);

        // Monitor connection state
        scheduledExecutor.scheduleAtFixedRate(this::checkConnectionState, 0, 5, TimeUnit.SECONDS);
    }

    private void checkConnectionState() {
        if (isShutdown.get()) {
            return;
        }

        try {
            ConnectivityState state = channel.getState(false);
            boolean wasConnected = isConnected.get();
            boolean nowConnected = state == ConnectivityState.READY || state == ConnectivityState.IDLE;

            if (wasConnected && !nowConnected) {
                LOG.warn("gRPC connection lost, state: {}", state);
                isConnected.set(false);
            } else if (!wasConnected && nowConnected) {
                LOG.info("gRPC connection restored");
                isConnected.set(true);

                // Restart task stream if needed
                startTaskStream();
            }
        } catch (Exception e) {
            LOG.warn("Error checking gRPC connection state", e);
        }
    }

    @Override
    public Duration getHeartbeatInterval() {
        return heartbeatInterval;
    }

    @Override
    public tech.kayys.silat.model.CommunicationType getCommunicationType() {
        return tech.kayys.silat.model.CommunicationType.GRPC;
    }

    @Override
    public Uni<Void> register(List<WorkflowExecutor> executors) {
        if (executors.isEmpty()) {
            return Uni.createFrom().voidItem();
        }

        WorkflowExecutor first = executors.get(0);
        RegisterExecutorRequest request = RegisterExecutorRequest.newBuilder()
                .setExecutorId(executorId)
                .setExecutorType(first.getExecutorType())
                .setCommunicationType(tech.kayys.silat.grpc.v1.CommunicationType.COMMUNICATION_TYPE_GRPC)
                .setEndpoint(java.net.InetAddress.getLoopbackAddress().getHostAddress())
                .setMaxConcurrentTasks(first.getMaxConcurrentTasks())
                .addAllSupportedNodeTypes(java.util.Arrays.asList(first.getSupportedNodeTypes()))
                .build();

        LOG.info("Registering executor {} via gRPC", executorId);

        return stub.registerExecutor(request)
                .onItem().invoke(resp -> LOG.info("Executor registered successfully with ID: {}", resp.getExecutorId()))
                .onFailure().invoke(error -> LOG.error("Failed to register executor {}", executorId, error))
                .replaceWithVoid();
    }

    @Override
    public Uni<Void> unregister() {
        UnregisterExecutorRequest request = UnregisterExecutorRequest.newBuilder()
                .setExecutorId(executorId)
                .build();

        LOG.info("Unregistering executor {} via gRPC", executorId);

        return stub.unregisterExecutor(request)
                .onItem().invoke(resp -> LOG.info("Executor unregistered successfully: {}", executorId))
                .onFailure().invoke(error -> LOG.error("Failed to unregister executor {}", executorId, error))
                .replaceWithVoid();
    }

    @Override
    public Multi<NodeExecutionTask> receiveTasks() {
        LOG.info("Setting up gRPC task stream for executor: {}", executorId);

        StreamTasksRequest request = StreamTasksRequest.newBuilder()
                .setExecutorId(executorId)
                .build();

        return stub.streamTasks(request)
                .onItem().transform(protoTask -> {
                    WorkflowRunId runId = WorkflowRunId.of(protoTask.getRunId());
                    NodeId nodeId = NodeId.of(protoTask.getNodeId());
                    int attempt = protoTask.getAttempt();
                    ExecutionToken token = new ExecutionToken(
                            protoTask.getExecutionToken(),
                            runId,
                            nodeId,
                            attempt,
                            Instant.now().plus(Duration.ofHours(1)));

                    return new NodeExecutionTask(
                            runId,
                            nodeId,
                            attempt,
                            token,
                            mapper.structToMap(protoTask.getContext()),
                            null // retryPolicy not provided in proto
                    );
                });
    }

    private void startTaskStream() {
        if (isShutdown.get()) {
            return;
        }
        LOG.info("Task stream setup initiated for executor: {}", executorId);
    }

    @Override
    public Uni<Void> sendResult(NodeExecutionResult result) {
        TaskResult protoResult = TaskResult.newBuilder()
                .setTaskId(result.getNodeId())
                .setRunId(result.runId().value())
                .setNodeId(result.getNodeId())
                .setAttempt(result.attempt())
                .setExecutionToken(result.executionToken().token())
                .setStatus(tech.kayys.silat.grpc.v1.TaskStatus.valueOf("TASK_STATUS_" + result.status().name()))
                .setOutput(mapper.mapToStruct(result.getUpdatedContext().getVariables()))
                .build();

        return stub.reportResults(Multi.createFrom().item(protoResult))
                .onItem().invoke(() -> LOG.debug("Result sent successfully for task: {}", result.getNodeId()))
                .onFailure().invoke(error -> LOG.error("Failed to send result for task: {}", result.getNodeId(), error))
                .replaceWithVoid();
    }

    @Override
    public Uni<Void> sendHeartbeat() {
        if (!isConnected.get()) {
            return Uni.createFrom().voidItem();
        }

        HeartbeatRequest request = HeartbeatRequest.newBuilder()
                .setExecutorId(executorId)
                .build();

        return stub.heartbeat(request)
                .onFailure().invoke(error -> LOG.warn("Heartbeat failed for executor: {}", executorId, error))
                .replaceWithVoid();
    }

    @PreDestroy
    public void cleanup() {
        LOG.info("Cleaning up gRPC transport for executor: {}", executorId);

        isShutdown.set(true);

        if (channel != null && !channel.isShutdown()) {
            try {
                channel.shutdown().awaitTermination(5, TimeUnit.SECONDS);
            } catch (InterruptedException e) {
                Thread.currentThread().interrupt();
                channel.shutdownNow();
            }
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/JwtClientInterceptor.java
Size: 1.1 KB | Modified: 2026-01-06 13:27:58
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import io.grpc.*;

/**
 * gRPC interceptor that adds a JWT token to the request metadata
 */
public class JwtClientInterceptor implements ClientInterceptor {

    private static final Metadata.Key<String> AUTHORIZATION_KEY = Metadata.Key.of("Authorization",
            Metadata.ASCII_STRING_MARSHALLER);

    private final String token;

    public JwtClientInterceptor(String token) {
        this.token = token;
    }

    @Override
    public <ReqT, RespT> ClientCall<ReqT, RespT> interceptCall(
            MethodDescriptor<ReqT, RespT> method,
            CallOptions callOptions,
            Channel next) {

        return new ForwardingClientCall.SimpleForwardingClientCall<ReqT, RespT>(next.newCall(method, callOptions)) {
            @Override
            public void start(Listener<RespT> responseListener, Metadata headers) {
                if (token != null && !token.isEmpty()) {
                    headers.put(AUTHORIZATION_KEY, "Bearer " + token);
                }
                super.start(responseListener, headers);
            }
        };
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/KafkaExecutorTransport.java
Size: 5.8 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.time.Duration;
import java.util.List;
import java.util.UUID;

import jakarta.annotation.PreDestroy;
import jakarta.enterprise.context.ApplicationScoped;

import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.eclipse.microprofile.reactive.messaging.Channel;
import org.eclipse.microprofile.reactive.messaging.Emitter;
import org.eclipse.microprofile.reactive.messaging.Incoming;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import io.smallrye.mutiny.operators.multi.processors.BroadcastProcessor;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;

/**
 * Kafka-based executor transport
 */
@ApplicationScoped
public class KafkaExecutorTransport implements ExecutorTransport {

    private static final Logger LOG = LoggerFactory.getLogger(KafkaExecutorTransport.class);

    private final String executorId;

    @ConfigProperty(name = "kafka.registration.timeout", defaultValue = "30s")
    Duration registrationTimeout;

    @ConfigProperty(name = "heartbeat.interval", defaultValue = "30s")
    Duration heartbeatInterval;

    // For task processing
    private final BroadcastProcessor<NodeExecutionTask> taskProcessor = BroadcastProcessor.create();

    // Kafka producers for different topics
    @Channel("execution-results")
    private Emitter<NodeExecutionResult> resultEmitter;

    @Channel("executor-heartbeats")
    private Emitter<ExecutorHeartbeat> heartbeatEmitter;

    @Channel("executor-registrations")
    private Emitter<ExecutorRegistration> registrationEmitter;

    @Channel("executor-unregistrations")
    private Emitter<ExecutorUnregistration> unregistrationEmitter;

    public KafkaExecutorTransport() {
        this.executorId = UUID.randomUUID().toString();
    }

    @Incoming("workflow-tasks")
    public void consumeTask(NodeExecutionTask task) {
        LOG.debug("Received task: {} from Kafka", task.nodeId());
        taskProcessor.onNext(task);
    }

    @Override
    public tech.kayys.silat.model.CommunicationType getCommunicationType() {
        return tech.kayys.silat.model.CommunicationType.KAFKA;
    }

    @Override
    public Uni<Void> register(List<WorkflowExecutor> executors) {
        LOG.info("Registering {} executors via Kafka", executors.size());

        List<String> executorTypes = executors.stream()
                .map(WorkflowExecutor::getExecutorType)
                .toList();

        ExecutorRegistration registration = new ExecutorRegistration(
                executorId,
                executorTypes,
                System.currentTimeMillis());

        return Uni.createFrom().completionStage(registrationEmitter.send(registration))
                .onItem().invoke(() -> LOG.info("Executor registration message sent: {}", executorId))
                .onFailure().invoke(e -> LOG.error("Failed to send registration message", e))
                .ifNoItem().after(registrationTimeout).fail()
                .replaceWithVoid();
    }

    @Override
    public Uni<Void> unregister() {
        LOG.info("Unregistering via Kafka for executor: {}", executorId);

        ExecutorUnregistration unregistration = new ExecutorUnregistration(
                executorId,
                System.currentTimeMillis());

        return Uni.createFrom().completionStage(unregistrationEmitter.send(unregistration))
                .onItem().invoke(() -> LOG.info("Executor unregistration message sent: {}", executorId))
                .onFailure().invoke(e -> LOG.error("Failed to send unregistration message", e))
                .replaceWithVoid();
    }

    @Override
    public Multi<NodeExecutionTask> receiveTasks() {
        LOG.info("Setting up Kafka task consumer");
        return taskProcessor;
    }

    @Override
    public Uni<Void> sendResult(NodeExecutionResult result) {
        return Uni.createFrom().emitter(emitter -> {
            try {
                // Send result to Kafka
                resultEmitter.send(result);

                LOG.debug("Result sent to Kafka for task: {}", result.getNodeId());
                emitter.complete(null);
            } catch (Exception e) {
                LOG.error("Failed to send result for task: {}", result.getNodeId(), e);
                emitter.fail(e);
            }
        });
    }

    @Override
    public Uni<Void> sendHeartbeat() {
        return Uni.createFrom().emitter(emitter -> {
            try {
                ExecutorHeartbeat heartbeat = new ExecutorHeartbeat(
                        executorId,
                        System.currentTimeMillis());

                // Send heartbeat to Kafka
                heartbeatEmitter.send(heartbeat);

                LOG.trace("Heartbeat sent to Kafka for executor: {}", executorId);
                emitter.complete(null);
            } catch (Exception e) {
                LOG.warn("Failed to send heartbeat for executor: {}", executorId, e);
                emitter.complete(null); // Don't fail for heartbeat issues
            }
        });
    }

    @Override
    public Duration getHeartbeatInterval() {
        return heartbeatInterval;
    }

    /**
     * Helper class for executor registration messages
     */
    public record ExecutorRegistration(String executorId, List<String> supportedTypes, long timestamp) {
    }

    /**
     * Helper class for executor unregistration messages
     */
    public record ExecutorUnregistration(String executorId, long timestamp) {
    }

    /**
     * Helper class for executor heartbeat messages
     */
    public record ExecutorHeartbeat(String executorId, long timestamp) {
    }

    @PreDestroy
    public void cleanup() {
        LOG.info("Cleaning up Kafka transport for executor: {}", executorId);

        // Close the task processor
        taskProcessor.onComplete();
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/LocalExecutorTransport.java
Size: 3.4 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.time.Duration;
import java.util.List;
import java.util.Map;

import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import io.vertx.mutiny.core.eventbus.EventBus;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;
import tech.kayys.silat.model.CommunicationType;
import tech.kayys.silat.model.ExecutorInfo;

/**
 * Transport implementation for local (same JVM) communication via Vert.x
 * EventBus
 */
@ApplicationScoped
public class LocalExecutorTransport implements ExecutorTransport {

    private static final Logger LOG = LoggerFactory.getLogger(LocalExecutorTransport.class);
    private static final String TOPIC_TASKS = "silat.tasks";
    private static final String TOPIC_RESULTS = "silat.results";
    private static final String TOPIC_REGISTER = "silat.executor.register";
    private static final String TOPIC_UNREGISTER = "silat.executor.unregister";
    private static final String TOPIC_HEARTBEAT = "silat.executor.heartbeat";

    private final java.util.Set<String> registeredExecutorIds = java.util.concurrent.ConcurrentHashMap.newKeySet();

    @Inject
    EventBus eventBus;

    @Override
    public tech.kayys.silat.model.CommunicationType getCommunicationType() {
        return tech.kayys.silat.model.CommunicationType.LOCAL;
    }

    @Override
    public Uni<Void> register(List<WorkflowExecutor> executors) {
        return Uni.createFrom().item(() -> {
            executors.forEach(executor -> {
                String executorId = "local-" + executor.getExecutorType();
                registeredExecutorIds.add(executorId);
                LOG.info("Registering local executor via EventBus: {} (id: {})", executor.getExecutorType(),
                        executorId);

                ExecutorInfo info = new ExecutorInfo(
                        executorId,
                        executor.getExecutorType(),
                        CommunicationType.LOCAL,
                        "local",
                        Duration.ofSeconds(30),
                        Map.of());

                eventBus.publish(TOPIC_REGISTER, io.vertx.core.json.JsonObject.mapFrom(info));
            });
            return null;
        });
    }

    @Override
    public Uni<Void> unregister() {
        return Uni.createFrom().item(() -> {
            LOG.info("Unregistering local executors via EventBus");
            registeredExecutorIds.clear();
            eventBus.publish(TOPIC_UNREGISTER, "all");
            return null;
        });
    }

    @Override
    public Multi<NodeExecutionTask> receiveTasks() {
        return eventBus.<io.vertx.core.json.JsonObject>consumer(TOPIC_TASKS)
                .toMulti()
                .map(msg -> msg.body().mapTo(NodeExecutionTask.class));
    }

    @Override
    public Uni<Void> sendResult(NodeExecutionResult result) {
        return Uni.createFrom().item(() -> {
            eventBus.publish(TOPIC_RESULTS, io.vertx.core.json.JsonObject.mapFrom(result));
            return null;
        });
    }

    @Override
    public Uni<Void> sendHeartbeat() {
        return Uni.createFrom().item(() -> {
            registeredExecutorIds.forEach(id -> eventBus.publish(TOPIC_HEARTBEAT, id));
            return null;
        });
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/PluginBasedExecutor.java
Size: 2.0 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;
import tech.kayys.silat.model.ErrorInfo;
import tech.kayys.silat.plugin.executor.ExecutorPlugin;

/**
 * Plugin-based executor that delegates task execution to loaded plugins
 * 
 * This executor discovers the appropriate plugin for each task and
 * delegates execution to it.
 */
@Executor(executorType = "plugin-based")
@ApplicationScoped
public class PluginBasedExecutor extends AbstractWorkflowExecutor {

    private static final Logger LOG = LoggerFactory.getLogger(PluginBasedExecutor.class);

    @Inject
    ExecutorPluginManager pluginManager;

    @Override
    public Uni<NodeExecutionResult> execute(NodeExecutionTask task) {
        String taskType = extractNodeType(task);
        LOG.debug("Finding plugin for task: {} (type: {})", task.nodeId(), taskType);

        // Find suitable plugin
        ExecutorPlugin plugin = pluginManager.findPlugin(task);

        if (plugin == null) {
            LOG.error("No plugin found for task: {} (type: {})", task.nodeId(), taskType);
            return Uni.createFrom().item(SimpleNodeExecutionResult.failure(
                    task.runId(),
                    task.nodeId(),
                    task.attempt(),
                    ErrorInfo.of(new IllegalStateException("No executor plugin found for task type: " + taskType)),
                    task.token()));
        }

        LOG.info("Executing task {} with plugin: {}", task.nodeId(), plugin.getMetadata().name());
        return plugin.execute(task);
    }

    @Override
    public boolean canHandle(NodeExecutionTask task) {
        boolean canHandle = pluginManager.hasPluginFor(task);
        LOG.debug("Can handle task {} (type: {}): {}", task.nodeId(), extractNodeType(task), canHandle);
        return canHandle;
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/SecurityConfig.java
Size: 1.9 KB | Modified: 2026-01-06 22:20:38
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

/**
 * Security configuration for the executor
 */
public record SecurityConfig(
        boolean mtlsEnabled,
        boolean jwtEnabled,
        String keyCertChainPath,
        String privateKeyPath,
        String trustCertCollectionPath,
        String jwtToken) {

    public static SecurityConfig disabled() {
        return builder()
                .mtlsEnabled(false)
                .jwtEnabled(false)
                .build();
    }

    public static Builder builder() {
        return new Builder();
    }

    public static class Builder {
        private boolean mtlsEnabled;
        private boolean jwtEnabled;
        private String keyCertChainPath;
        private String privateKeyPath;
        private String trustCertCollectionPath;
        private String jwtToken;

        public Builder mtlsEnabled(boolean mtlsEnabled) {
            this.mtlsEnabled = mtlsEnabled;
            return this;
        }

        public Builder jwtEnabled(boolean jwtEnabled) {
            this.jwtEnabled = jwtEnabled;
            return this;
        }

        public Builder keyCertChainPath(String keyCertChainPath) {
            this.keyCertChainPath = keyCertChainPath;
            return this;
        }

        public Builder privateKeyPath(String privateKeyPath) {
            this.privateKeyPath = privateKeyPath;
            return this;
        }

        public Builder trustCertCollectionPath(String trustCertCollectionPath) {
            this.trustCertCollectionPath = trustCertCollectionPath;
            return this;
        }

        public Builder jwtToken(String jwtToken) {
            this.jwtToken = jwtToken;
            return this;
        }

        public SecurityConfig build() {
            return new SecurityConfig(mtlsEnabled, jwtEnabled, keyCertChainPath, privateKeyPath, trustCertCollectionPath, jwtToken);
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/SimpleNodeExecutionResult.java
Size: 3.8 KB | Modified: 2026-01-03 19:00:42
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import java.time.Duration;
import java.time.Instant;
import java.util.Map;

import tech.kayys.silat.execution.ExecutionContext;
import tech.kayys.silat.execution.ExecutionError;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionStatus;
import tech.kayys.silat.model.ErrorInfo;
import tech.kayys.silat.model.ExecutionToken;
import tech.kayys.silat.model.NodeId;
import tech.kayys.silat.model.WaitInfo;
import tech.kayys.silat.model.WorkflowRunId;

/**
 * Simple implementation of NodeExecutionResult for executor SDK
 */
public record SimpleNodeExecutionResult(
        WorkflowRunId runId,
        NodeId nodeId,
        int attempt,
        NodeExecutionStatus status,
        Map<String, Object> output,
        ErrorInfo error,
        ExecutionToken executionToken,
        Instant executedAt,
        Duration duration,
        ExecutionContext updatedContext,
        WaitInfo waitInfo,
        Map<String, Object> metadata) implements NodeExecutionResult {

    @Override
    public NodeExecutionStatus getStatus() {
        return status;
    }

    @Override
    public String getNodeId() {
        return nodeId != null ? nodeId.value() : null;
    }

    @Override
    public Instant getExecutedAt() {
        return executedAt;
    }

    @Override
    public Duration getDuration() {
        return duration;
    }

    @Override
    public ExecutionContext getUpdatedContext() {
        return updatedContext;
    }

    @Override
    public ExecutionError getError() {
        if (error == null) {
            return null;
        }
        return new ExecutionError() {
            @Override
            public String getCode() {
                return error.code();
            }

            @Override
            public Category getCategory() {
                return Category.SYSTEM;
            }

            @Override
            public String getMessage() {
                return error.message();
            }

            @Override
            public boolean isRetriable() {
                return false;
            }

            @Override
            public String getCompensationHint() {
                return null;
            }

            @Override
            public Map<String, Object> getDetails() {
                return error.context();
            }
        };
    }

    @Override
    public WaitInfo getWaitInfo() {
        return waitInfo;
    }

    @Override
    public Map<String, Object> getMetadata() {
        return metadata;
    }

    /**
     * Create a failure result
     */
    public static NodeExecutionResult failure(
            WorkflowRunId runId,
            NodeId nodeId,
            int attempt,
            ErrorInfo error,
            ExecutionToken token) {
        return new SimpleNodeExecutionResult(
                runId,
                nodeId,
                attempt,
                NodeExecutionStatus.FAILED,
                Map.of(),
                error,
                token,
                Instant.now(),
                Duration.ZERO,
                null,
                null,
                Map.of());
    }

    /**
     * Create a success result
     */
    public static NodeExecutionResult success(
            WorkflowRunId runId,
            NodeId nodeId,
            int attempt,
            Map<String, Object> output,
            ExecutionToken token,
            Duration duration) {
        return new SimpleNodeExecutionResult(
                runId,
                nodeId,
                attempt,
                NodeExecutionStatus.COMPLETED,
                output,
                null,
                token,
                Instant.now(),
                duration,
                null,
                null,
                Map.of());
    }
}

================================================================================

================================================================================
tech/kayys/silat/sdk/executor/WorkflowExecutor.java
Size: 2.1 KB | Modified: 2026-01-03 18:52:04
--------------------------------------------------------------------------------
package tech.kayys.silat.sdk.executor;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;

/**
 * Base interface for all workflow executors
 */
public interface WorkflowExecutor {

    /**
     * Execute a node task
     *
     * @param task The task to execute
     * @return Result of execution
     */
    Uni<NodeExecutionResult> execute(NodeExecutionTask task);

    /**
     * Get executor type
     */
    String getExecutorType();

    /**
     * Validate if this executor can handle the task
     */
    default boolean canHandle(NodeExecutionTask task) {
        return true;
    }

    /**
     * Get the maximum number of concurrent tasks this executor can handle
     */
    default int getMaxConcurrentTasks() {
        return Integer.MAX_VALUE; // Unlimited by default
    }

    /**
     * Called before execution starts
     */
    default Uni<Void> beforeExecute(NodeExecutionTask task) {
        return Uni.createFrom().voidItem();
    }

    /**
     * Called after execution completes (success or failure)
     */
    default Uni<Void> afterExecute(
            NodeExecutionTask task,
            NodeExecutionResult result) {
        return Uni.createFrom().voidItem();
    }

    /**
     * Called when execution fails with an exception
     */
    default Uni<Void> onError(NodeExecutionTask task, Throwable error) {
        return Uni.createFrom().voidItem();
    }

    /**
     * Get supported node types that this executor can handle
     */
    default String[] getSupportedNodeTypes() {
        return new String[0]; // Empty array means all types supported
    }

    /**
     * Check if the executor is ready to accept new tasks
     */
    default boolean isReady() {
        return true;
    }

    /**
     * Initialize the executor (called during registration)
     */
    default Uni<Void> initialize() {
        return Uni.createFrom().voidItem();
    }

    /**
     * Cleanup the executor (called during unregistration)
     */
    default Uni<Void> cleanup() {
        return Uni.createFrom().voidItem();
    }
}
================================================================================

