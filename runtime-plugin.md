Coto Output
Generated: 2026-01-17 09:26:44
Files: 8 | Directories: 9 | Total Size: 21.6 KB


================================================================================
tech/kayys/silat/runtime/DbInitializer.java
Size: 1.8 KB | Modified: 2026-01-14 11:57:13
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime;

import io.quarkus.runtime.StartupEvent;
import io.vertx.mutiny.pgclient.PgPool;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.event.Observes;
import jakarta.inject.Inject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

/**
 * Automatically initializes the database schema on startup.
 * Useful when not using Flyway or Liquibase for non-entity tables.
 */
@ApplicationScoped
public class DbInitializer {
    private static final Logger LOGGER = LoggerFactory.getLogger(DbInitializer.class);

    @Inject
    PgPool client;

    void onStart(@Observes StartupEvent ev) {
        LOGGER.info("Initializing database schema for workflow_definitions...");

        String sql = """
                CREATE TABLE IF NOT EXISTS workflow_definitions (
                    definition_id VARCHAR(128) PRIMARY KEY,
                    tenant_id VARCHAR(64) NOT NULL,
                    name VARCHAR(255) NOT NULL,
                    version VARCHAR(32) NOT NULL,
                    description TEXT,
                    definition_json JSONB NOT NULL,
                    is_active BOOLEAN DEFAULT true,
                    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                    created_by VARCHAR(128),
                    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
                    updated_by VARCHAR(128),
                    metadata JSONB,
                    CONSTRAINT uk_workflow_def_tenant_name_version UNIQUE (tenant_id, name, version)
                );
                """;

        client.query(sql).execute()
                .subscribe().with(
                        result -> LOGGER.info("Database schema 'workflow_definitions' initialized successfully"),
                        error -> LOGGER.error("Failed to initialize database schema", error));
    }
}

================================================================================

================================================================================
tech/kayys/silat/runtime/grpc/ExecutorServiceImpl.java
Size: 3.7 KB | Modified: 2026-01-15 14:57:23
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.grpc;

import com.google.protobuf.Empty;
import io.quarkus.grpc.GrpcService;
import io.smallrye.mutiny.Multi;
import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.kayys.silat.grpc.v1.*;
import tech.kayys.silat.model.CommunicationType;
import tech.kayys.silat.grpc.CommunicationTypeConverter;
import tech.kayys.silat.model.ExecutorInfo;
import tech.kayys.silat.registry.ExecutorRegistryService;

import java.time.Duration;
import java.util.Map;

@GrpcService
public class ExecutorServiceImpl extends MutinyExecutorServiceGrpc.ExecutorServiceImplBase {

    private static final Logger LOG = LoggerFactory.getLogger(ExecutorServiceImpl.class);

    @Inject
    ExecutorRegistryService executorRegistry;

    public ExecutorServiceImpl() {
        System.out.println("ExecutorServiceImpl initialized!");
        LOG.info("ExecutorServiceImpl initialized!");
    }

    @Override
    public Uni<ExecutorRegistration> registerExecutor(RegisterExecutorRequest request) {
        System.out.println("ExecutorServiceImpl: registerExecutor called for " + request.getExecutorId());
        LOG.info("Received registration request from executor: {}", request.getExecutorId());

        ExecutorInfo executorInfo = new ExecutorInfo(
                request.getExecutorId(),
                request.getExecutorType(),
                mapCommunicationType(request.getCommunicationType()),
                request.getEndpoint(),
                Duration.ofHours(24), // TODO: map from request if available, or use default
                Map.of() // Metadata
        );

        return executorRegistry.registerExecutor(executorInfo)
                .map(v -> ExecutorRegistration.newBuilder()
                        .setExecutorId(request.getExecutorId())
                        .setRegisteredAt(com.google.protobuf.Timestamp.newBuilder()
                                .setSeconds(System.currentTimeMillis() / 1000)
                                .setNanos((int) ((System.currentTimeMillis() % 1000) * 1000000))
                                .build())
                        .build());
    }

    @Override
    public Uni<Empty> unregisterExecutor(UnregisterExecutorRequest request) {
        LOG.info("Received unregistration request from executor: {}", request.getExecutorId());
        return executorRegistry.unregisterExecutor(request.getExecutorId())
                .map(v -> Empty.getDefaultInstance());
    }

    @Override
    public Uni<Empty> heartbeat(HeartbeatRequest request) {
        LOG.trace("Received heartbeat from executor: {}", request.getExecutorId());
        return executorRegistry.heartbeat(request.getExecutorId())
                .map(v -> Empty.getDefaultInstance());
    }

    @Override
    public Multi<ExecutionTask> streamTasks(StreamTasksRequest request) {
        LOG.info("Executor {} requested task stream", request.getExecutorId());

        // For now, return an empty stream that won't cause null pointer exceptions
        // TODO: Implement actual task streaming logic
        return Multi.createFrom().empty();
    }

    @Override
    public Uni<Empty> reportResults(Multi<TaskResult> request) {
        return request.onItem().invoke(result -> {
            LOG.info("Received result for task: {} from executor: {}", result.getTaskId(), "UNKNOWN");
            // TODO: Process result
        }).collect().last().map(v -> Empty.getDefaultInstance());
    }

    @Override
    public Multi<EngineMessage> executeStream(Multi<ExecutorMessage> request) {
        return Multi.createFrom().empty();
    }

    private CommunicationType mapCommunicationType(tech.kayys.silat.grpc.v1.CommunicationType grpcType) {
        return CommunicationTypeConverter.fromGrpc(grpcType);
    }
}

================================================================================

================================================================================
tech/kayys/silat/runtime/resource/CallbackResource.java
Size: 964 B | Modified: 2026-01-14 07:26:15
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.resource;

import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import tech.kayys.silat.api.engine.WorkflowRunManager;
import tech.kayys.silat.execution.ExternalSignal;
import tech.kayys.silat.model.WorkflowRunId;

@Path("/api/v1/callbacks")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class CallbackResource {

    @Inject
    WorkflowRunManager runManager;

    @POST
    @Path("/{runId}/signal")
    public Uni<Void> signal(
            @PathParam("runId") String runId,
            @QueryParam("token") String callbackToken,
            ExternalSignal signal) {
        return runManager.onExternalSignal(WorkflowRunId.of(runId), signal, callbackToken);
    }
}

================================================================================

================================================================================
tech/kayys/silat/runtime/resource/ExecutorRegistryResource.java
Size: 4.7 KB | Modified: 2026-01-14 12:34:58
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.resource;

import java.util.List;
import java.util.Map;
import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.DELETE;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.PUT;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import tech.kayys.silat.model.ExecutorHealthInfo;
import tech.kayys.silat.model.ExecutorInfo;
import tech.kayys.silat.model.NodeId;
import tech.kayys.silat.registry.ExecutorRegistryService;
import tech.kayys.silat.registry.ExecutorStatistics;

@Path("/api/v1/executors")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class ExecutorRegistryResource {

    @Inject
    ExecutorRegistryService executorRegistryService;

    @POST
    public Uni<Void> registerExecutor(ExecutorInfo executor) {
        return executorRegistryService.registerExecutor(executor);
    }

    @DELETE
    @Path("/{executorId}")
    public Uni<Void> unregisterExecutor(@PathParam("executorId") String executorId) {
        return executorRegistryService.unregisterExecutor(executorId);
    }

    @POST
    @Path("/{executorId}/heartbeat")
    public Uni<Void> heartbeat(@PathParam("executorId") String executorId) {
        return executorRegistryService.heartbeat(executorId);
    }

    @GET
    @Path("/{executorId}")
    public Uni<ExecutorInfo> getExecutorById(@PathParam("executorId") String executorId) {
        return executorRegistryService.getExecutorById(executorId)
                .map(optional -> optional.orElse(null)); // Return null if not found
    }

    @GET
    public Uni<List<ExecutorInfo>> getAllExecutors(
            @QueryParam("healthyOnly") Boolean healthyOnly,
            @QueryParam("type") String type,
            @QueryParam("communicationType") String communicationType) {

        Uni<List<ExecutorInfo>> result;

        if (healthyOnly != null && healthyOnly) {
            result = executorRegistryService.getHealthyExecutors();
        } else {
            result = executorRegistryService.getAllExecutors();
        }

        // Apply filters if specified
        if (type != null) {
            result = result.map(executors -> executors.stream()
                    .filter(executor -> executor.executorType().equals(type))
                    .toList());
        }

        if (communicationType != null) {
            result = result.map(executors -> executors.stream()
                    .filter(executor -> executor.communicationType().toString().equalsIgnoreCase(communicationType))
                    .toList());
        }

        return result;
    }

    @GET
    @Path("/healthy")
    public Uni<List<ExecutorInfo>> getHealthyExecutors() {
        return executorRegistryService.getHealthyExecutors();
    }

    @GET
    @Path("/count")
    public Uni<Integer> getExecutorCount() {
        return executorRegistryService.getExecutorCount();
    }

    @GET
    @Path("/statistics")
    public Uni<ExecutorStatistics> getStatistics() {
        return executorRegistryService.getStatistics();
    }

    @GET
    @Path("/type/{type}")
    public Uni<List<ExecutorInfo>> getExecutorsByType(@PathParam("type") String type) {
        return executorRegistryService.getExecutorsByType(type);
    }

    @GET
    @Path("/communication-type/{communicationType}")
    public Uni<List<ExecutorInfo>> getExecutorsByCommunicationType(
            @PathParam("communicationType") String communicationType) {
        return executorRegistryService.getExecutorsByCommunicationType(
                tech.kayys.silat.model.CommunicationType.valueOf(communicationType.toUpperCase()));
    }

    @GET
    @Path("/health/{executorId}")
    public Uni<ExecutorHealthInfo> getHealthInfo(@PathParam("executorId") String executorId) {
        return executorRegistryService.getHealthInfo(executorId)
                .map(optional -> optional.orElse(null)); // Return null if not found
    }

    @GET
    @Path("/healthy/{executorId}")
    public Uni<Boolean> isHealthy(@PathParam("executorId") String executorId) {
        return executorRegistryService.isHealthy(executorId);
    }

    @PUT
    @Path("/{executorId}/metadata")
    public Uni<Void> updateExecutorMetadata(@PathParam("executorId") String executorId,
            Map<String, String> metadata) {
        return executorRegistryService.updateExecutorMetadata(executorId, metadata);
    }

    @POST
    @Path("/select/{nodeId}")
    public Uni<ExecutorInfo> getExecutorForNode(@PathParam("nodeId") String nodeId) {
        return executorRegistryService.getExecutorForNode(new NodeId(nodeId))
                .map(optional -> optional.orElse(null)); // Return null if not found
    }
}
================================================================================

================================================================================
tech/kayys/silat/runtime/resource/PluginResource.java
Size: 2.2 KB | Modified: 2026-01-14 14:30:56
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.resource;

import java.util.List;

import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.core.MediaType;
import tech.kayys.silat.plugin.Plugin;
import tech.kayys.silat.plugin.PluginService;

@Path("/api/v1/plugins")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class PluginResource {

    @Inject
    PluginService pluginService;

    @GET
    public Uni<List<Plugin>> getAllPlugins() {
        return Uni.createFrom().item(pluginService.getAllPlugins());
    }

    @GET
    @Path("/{pluginId}")
    public Uni<Plugin> getPlugin(@PathParam("pluginId") String pluginId) {
        return Uni.createFrom().item(pluginService.getPlugin(pluginId).orElse(null));
    }

    @POST
    @Path("/{pluginId}/start")
    public Uni<Void> startPlugin(@PathParam("pluginId") String pluginId) {
        return pluginService.startPlugin(pluginId);
    }

    @POST
    @Path("/{pluginId}/stop")
    public Uni<Void> stopPlugin(@PathParam("pluginId") String pluginId) {
        return pluginService.stopPlugin(pluginId);
    }

    @GET
    @Path("/types/{pluginType}")
    public Uni<List<Plugin>> getPluginsByType(@PathParam("pluginType") String pluginType) {
        // This would require reflection to determine plugin types
        // For now, we'll return all plugins
        return Uni.createFrom().item(pluginService.getAllPlugins());
    }

    @GET
    @Path("/status")
    public Uni<List<PluginStatusInfo>> getPluginStatuses() {
        List<Plugin> plugins = pluginService.getAllPlugins();
        List<PluginStatusInfo> statuses = plugins.stream()
                .map(plugin -> new PluginStatusInfo(
                        plugin.getMetadata().id(),
                        plugin.getMetadata().name(),
                        plugin.getMetadata().version(),
                        "ACTIVE" // Simplified status
                ))
                .toList();
        return Uni.createFrom().item(statuses);
    }

    public record PluginStatusInfo(String id, String name, String version, String status) {
    }
}
================================================================================

================================================================================
tech/kayys/silat/runtime/resource/WorkflowDefinitionResource.java
Size: 2.2 KB | Modified: 2026-01-14 11:57:28
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.resource;

import java.util.List;

import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.DELETE;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.PUT;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import tech.kayys.silat.dto.CreateWorkflowDefinitionRequest;
import tech.kayys.silat.dto.UpdateWorkflowDefinitionRequest;
import tech.kayys.silat.model.TenantId;
import tech.kayys.silat.model.WorkflowDefinition;
import tech.kayys.silat.model.WorkflowDefinitionId;
import tech.kayys.silat.runtime.workflow.RuntimeWorkflowDefinitionService;
import tech.kayys.silat.security.TenantSecurityContext;

@Path("/api/v1/workflow-definitions")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class WorkflowDefinitionResource {

    @Inject
    RuntimeWorkflowDefinitionService service;

    @Inject
    TenantSecurityContext securityContext;

    @POST
    public Uni<WorkflowDefinition> create(CreateWorkflowDefinitionRequest request) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return service.create(request, tenantId);
    }

    @GET
    @Path("/{id}")
    public Uni<WorkflowDefinition> get(@PathParam("id") String id) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return service.get(new WorkflowDefinitionId(id), tenantId);
    }

    @GET
    public Uni<List<WorkflowDefinition>> list(@QueryParam("activeOnly") boolean activeOnly) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return service.list(tenantId, activeOnly);
    }

    @PUT
    @Path("/{id}")
    public Uni<WorkflowDefinition> update(@PathParam("id") String id, UpdateWorkflowDefinitionRequest request) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return service.update(new WorkflowDefinitionId(id), request, tenantId);
    }

    @DELETE
    @Path("/{id}")
    public Uni<Void> delete(@PathParam("id") String id) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return service.delete(new WorkflowDefinitionId(id), tenantId);
    }
}

================================================================================

================================================================================
tech/kayys/silat/runtime/resource/WorkflowRunResource.java
Size: 4.1 KB | Modified: 2026-01-14 07:26:14
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.resource;

import java.util.List;
import java.util.Map;

import io.smallrye.mutiny.Uni;
import jakarta.inject.Inject;
import jakarta.ws.rs.Consumes;
import jakarta.ws.rs.GET;
import jakarta.ws.rs.POST;
import jakarta.ws.rs.Path;
import jakarta.ws.rs.PathParam;
import jakarta.ws.rs.Produces;
import jakarta.ws.rs.QueryParam;
import jakarta.ws.rs.core.MediaType;
import tech.kayys.silat.api.engine.WorkflowRunManager;
import tech.kayys.silat.execution.ExecutionHistory;
import tech.kayys.silat.model.CreateRunRequest;
import tech.kayys.silat.model.RunStatus;
import tech.kayys.silat.model.TenantId;
import tech.kayys.silat.model.WorkflowDefinitionId;
import tech.kayys.silat.model.WorkflowRun;
import tech.kayys.silat.model.WorkflowRunId;
import tech.kayys.silat.model.WorkflowRunSnapshot;
import tech.kayys.silat.security.TenantSecurityContext;

@Path("/api/v1/workflow-runs")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.APPLICATION_JSON)
public class WorkflowRunResource {

    @Inject
    WorkflowRunManager runManager;

    @Inject
    TenantSecurityContext securityContext;

    @POST
    public Uni<WorkflowRun> create(CreateRunRequest request) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return runManager.createRun(request, tenantId)
                .flatMap(run -> {
                    if (request.isAutoStart()) {
                        return runManager.startRun(run.getId(), tenantId);
                    }
                    return Uni.createFrom().item(run);
                });
    }

    @GET
    @Path("/{id}")
    public Uni<WorkflowRun> get(@PathParam("id") String id) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return runManager.getRun(WorkflowRunId.of(id), tenantId);
    }

    @GET
    @Path("/{id}/snapshot")
    public Uni<WorkflowRunSnapshot> getSnapshot(@PathParam("id") String id) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return runManager.getSnapshot(WorkflowRunId.of(id), tenantId);
    }

    @GET
    @Path("/{id}/history")
    public Uni<ExecutionHistory> getHistory(@PathParam("id") String id) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return runManager.getExecutionHistory(WorkflowRunId.of(id), tenantId);
    }

    @POST
    @Path("/{id}/start")
    public Uni<WorkflowRun> start(@PathParam("id") String id) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return runManager.startRun(WorkflowRunId.of(id), tenantId);
    }

    @POST
    @Path("/{id}/suspend")
    public Uni<WorkflowRun> suspend(@PathParam("id") String id, Map<String, Object> params) {
        TenantId tenantId = securityContext.getCurrentTenant();
        String reason = (String) params.getOrDefault("reason", "Manual suspension");
        // nodeId is optional for manual suspension
        return runManager.suspendRun(WorkflowRunId.of(id), tenantId, reason, null);
    }

    @POST
    @Path("/{id}/resume")
    public Uni<WorkflowRun> resume(@PathParam("id") String id, Map<String, Object> resumeData) {
        TenantId tenantId = securityContext.getCurrentTenant();
        return runManager.resumeRun(WorkflowRunId.of(id), tenantId, resumeData);
    }

    @POST
    @Path("/{id}/cancel")
    public Uni<Void> cancel(@PathParam("id") String id, Map<String, Object> params) {
        TenantId tenantId = securityContext.getCurrentTenant();
        String reason = (String) params.getOrDefault("reason", "Manual cancellation");
        return runManager.cancelRun(WorkflowRunId.of(id), tenantId, reason);
    }

    @GET
    public Uni<List<WorkflowRun>> query(
            @QueryParam("definitionId") String definitionId,
            @QueryParam("status") RunStatus status,
            @QueryParam("page") @jakarta.ws.rs.DefaultValue("0") int page,
            @QueryParam("size") @jakarta.ws.rs.DefaultValue("10") int size) {
        TenantId tenantId = securityContext.getCurrentTenant();
        WorkflowDefinitionId wfDefId = definitionId != null ? new WorkflowDefinitionId(definitionId) : null;
        return runManager.queryRuns(tenantId, wfDefId, status, page, size);
    }
}

================================================================================

================================================================================
tech/kayys/silat/runtime/workflow/RuntimeWorkflowDefinitionService.java
Size: 1.8 KB | Modified: 2026-01-13 07:11:22
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.workflow;

import io.smallrye.mutiny.Uni;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import tech.kayys.silat.dto.CreateWorkflowDefinitionRequest;
import tech.kayys.silat.dto.UpdateWorkflowDefinitionRequest;
import tech.kayys.silat.model.TenantId;
import tech.kayys.silat.model.WorkflowDefinition;
import tech.kayys.silat.model.WorkflowDefinitionId;
import java.util.List;

/**
 * Runtime service that acts as an adapter between API layer and engine layer
 */
@ApplicationScoped
public class RuntimeWorkflowDefinitionService {

    @Inject
    tech.kayys.silat.api.workflow.WorkflowDefinitionService engineService;

    public Uni<WorkflowDefinition> create(
            CreateWorkflowDefinitionRequest request,
            TenantId tenantId) {
        // The engine service now accepts DTOs directly since it implements the API
        // interface
        return engineService.create(request, tenantId);
    }

    public Uni<WorkflowDefinition> get(
            WorkflowDefinitionId id,
            TenantId tenantId) {
        return engineService.get(id, tenantId);
    }

    public Uni<List<WorkflowDefinition>> list(
            TenantId tenantId,
            boolean activeOnly) {
        return engineService.list(tenantId, activeOnly);
    }

    public Uni<WorkflowDefinition> update(
            WorkflowDefinitionId id,
            UpdateWorkflowDefinitionRequest request,
            TenantId tenantId) {
        // The engine service now accepts DTOs directly since it implements the API
        // interface
        return engineService.update(id, request, tenantId);
    }

    public Uni<Void> delete(
            WorkflowDefinitionId id,
            TenantId tenantId) {
        return engineService.delete(id, tenantId);
    }
}
================================================================================


Coto Output
Generated: 2026-01-17 15:51:27
Files: 23 | Directories: 13 | Total Size: 48.8 KB


================================================================================
tech/kayys/silat/plugin/EventBus.java
Size: 848 B | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

import java.util.function.Consumer;

/**
 * Event bus for plugin event communication
 * 
 * Plugins can publish and subscribe to events.
 */
public interface EventBus {
    
    /**
     * Publish an event
     * 
     * @param event the event to publish
     */
    void publish(PluginEvent event);
    
    /**
     * Subscribe to events of a specific type
     * 
     * @param eventType the event type
     * @param handler the event handler
     * @param <T> the event type
     * @return a subscription that can be used to unsubscribe
     */
    <T extends PluginEvent> Subscription subscribe(Class<T> eventType, Consumer<T> handler);
    
    /**
     * Subscription handle
     */
    interface Subscription {
        /**
         * Unsubscribe from events
         */
        void unsubscribe();
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/Plugin.java
Size: 1.4 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

/**
 * Base interface for all Silat plugins
 * 
 * Plugins must implement this interface and provide a no-arg constructor.
 * The plugin lifecycle is:
 * 1. Constructor called
 * 2. initialize(PluginContext) called
 * 3. start() called
 * 4. Plugin is active
 * 5. stop() called
 * 6. Plugin is unloaded
 */
public interface Plugin {
    
    /**
     * Initialize the plugin with the provided context
     * 
     * This method is called once after the plugin is loaded.
     * Use this to set up any resources needed by the plugin.
     * 
     * @param context the plugin context
     * @throws PluginException if initialization fails
     */
    void initialize(PluginContext context) throws PluginException;
    
    /**
     * Start the plugin
     * 
     * This method is called after initialization.
     * The plugin should start any background tasks or services here.
     * 
     * @throws PluginException if start fails
     */
    void start() throws PluginException;
    
    /**
     * Stop the plugin
     * 
     * This method is called when the plugin is being unloaded.
     * The plugin should clean up any resources and stop any background tasks.
     * 
     * @throws PluginException if stop fails
     */
    void stop() throws PluginException;
    
    /**
     * Get the plugin metadata
     * 
     * @return the plugin metadata
     */
    PluginMetadata getMetadata();
}

================================================================================

================================================================================
tech/kayys/silat/plugin/PluginContext.java
Size: 1.3 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

import org.slf4j.Logger;
import java.util.Map;
import java.util.Optional;

/**
 * Plugin context providing access to engine services and configuration
 */
public interface PluginContext {
    
    /**
     * Get the plugin metadata
     */
    PluginMetadata getMetadata();
    
    /**
     * Get a logger instance for this plugin
     */
    Logger getLogger();
    
    /**
     * Get a configuration property
     * 
     * @param key the property key
     * @return the property value if present
     */
    Optional<String> getProperty(String key);
    
    /**
     * Get a configuration property with a default value
     * 
     * @param key the property key
     * @param defaultValue the default value if property is not found
     * @return the property value or default
     */
    String getProperty(String key, String defaultValue);
    
    /**
     * Get all configuration properties
     */
    Map<String, String> getAllProperties();
    
    /**
     * Get the service registry for inter-plugin communication
     */
    ServiceRegistry getServiceRegistry();
    
    /**
     * Get the event bus for publishing/subscribing to events
     */
    EventBus getEventBus();
    
    /**
     * Get the plugin data directory for storing plugin-specific data
     */
    String getDataDirectory();
}

================================================================================

================================================================================
tech/kayys/silat/plugin/PluginEvent.java
Size: 949 B | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

import java.time.Instant;
import java.util.Map;

/**
 * Base class for plugin events
 */
public abstract class PluginEvent {
    
    private final String eventId;
    private final String sourcePluginId;
    private final Instant timestamp;
    private final Map<String, Object> metadata;
    
    protected PluginEvent(String sourcePluginId, Map<String, Object> metadata) {
        this.eventId = java.util.UUID.randomUUID().toString();
        this.sourcePluginId = sourcePluginId;
        this.timestamp = Instant.now();
        this.metadata = metadata != null ? Map.copyOf(metadata) : Map.of();
    }
    
    public String getEventId() {
        return eventId;
    }
    
    public String getSourcePluginId() {
        return sourcePluginId;
    }
    
    public Instant getTimestamp() {
        return timestamp;
    }
    
    public Map<String, Object> getMetadata() {
        return metadata;
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/PluginException.java
Size: 660 B | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

/**
 * Exception thrown by plugin operations
 */
public class PluginException extends Exception {
    
    private final String pluginId;
    
    public PluginException(String pluginId, String message) {
        super(message);
        this.pluginId = pluginId;
    }
    
    public PluginException(String pluginId, String message, Throwable cause) {
        super(message, cause);
        this.pluginId = pluginId;
    }
    
    public PluginException(String pluginId, Throwable cause) {
        super(cause);
        this.pluginId = pluginId;
    }
    
    public String getPluginId() {
        return pluginId;
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/PluginMetadata.java
Size: 1.5 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

import java.util.List;
import java.util.Map;

/**
 * Plugin metadata containing information about a plugin
 */
public record PluginMetadata(
    String id,
    String name,
    String version,
    String author,
    String description,
    List<PluginDependency> dependencies,
    Map<String, String> properties
) {
    public PluginMetadata {
        if (id == null || id.isBlank()) {
            throw new IllegalArgumentException("Plugin ID cannot be null or blank");
        }
        if (name == null || name.isBlank()) {
            throw new IllegalArgumentException("Plugin name cannot be null or blank");
        }
        if (version == null || version.isBlank()) {
            throw new IllegalArgumentException("Plugin version cannot be null or blank");
        }
        dependencies = dependencies != null ? List.copyOf(dependencies) : List.of();
        properties = properties != null ? Map.copyOf(properties) : Map.of();
    }

    /**
     * Plugin dependency information
     */
    public record PluginDependency(
        String pluginId,
        String versionConstraint
    ) {
        public PluginDependency {
            if (pluginId == null || pluginId.isBlank()) {
                throw new IllegalArgumentException("Plugin dependency ID cannot be null or blank");
            }
            if (versionConstraint == null || versionConstraint.isBlank()) {
                throw new IllegalArgumentException("Version constraint cannot be null or blank");
            }
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/PluginMetadataBuilder.java
Size: 1.6 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

/**
 * Builder for PluginMetadata
 */
public class PluginMetadataBuilder {
    private String id;
    private String name;
    private String version;
    private String author;
    private String description;
    private List<PluginMetadata.PluginDependency> dependencies = new ArrayList<>();
    private Map<String, String> properties = new HashMap<>();

    public static PluginMetadataBuilder builder() {
        return new PluginMetadataBuilder();
    }

    public PluginMetadataBuilder id(String id) {
        this.id = id;
        return this;
    }

    public PluginMetadataBuilder name(String name) {
        this.name = name;
        return this;
    }

    public PluginMetadataBuilder version(String version) {
        this.version = version;
        return this;
    }

    public PluginMetadataBuilder author(String author) {
        this.author = author;
        return this;
    }

    public PluginMetadataBuilder description(String description) {
        this.description = description;
        return this;
    }

    public PluginMetadataBuilder addDependency(String pluginId, String versionConstraint) {
        this.dependencies.add(new PluginMetadata.PluginDependency(pluginId, versionConstraint));
        return this;
    }

    public PluginMetadataBuilder property(String key, String value) {
        this.properties.put(key, value);
        return this;
    }

    public PluginMetadata build() {
        return new PluginMetadata(id, name, version, author, description, dependencies, properties);
    }
}
================================================================================

================================================================================
tech/kayys/silat/plugin/PluginService.java
Size: 1.5 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

import java.nio.file.Path;
import java.util.List;
import java.util.Optional;

import io.smallrye.mutiny.Uni;

/**
 * Unified plugin service that combines all plugin management functionality
 * This provides a single entry point for plugin operations in the engine
 */
public interface PluginService extends ServiceRegistry, EventBus {

    /**
     * Load a plugin from a JAR file
     */
    Uni<Plugin> loadPlugin(Path pluginJar);

    /**
     * Register a plugin instance directly (programmatic registration)
     */
    Uni<Void> registerPlugin(Plugin plugin);

    /**
     * Start a plugin
     */
    Uni<Void> startPlugin(String pluginId);

    /**
     * Stop a plugin
     */
    Uni<Void> stopPlugin(String pluginId);

    /**
     * Unload a plugin
     */
    Uni<Void> unloadPlugin(String pluginId);

    /**
     * Reload a plugin (hot-reload)
     */
    Uni<Plugin> reloadPlugin(String pluginId, Path pluginJar);

    /**
     * Get a plugin by ID
     */
    Optional<Plugin> getPlugin(String pluginId);

    /**
     * Get all loaded plugins
     */
    List<Plugin> getAllPlugins();

    /**
     * Get plugins by type
     */
    <T extends Plugin> List<T> getPluginsByType(Class<T> pluginType);

    /**
     * Discover and load all plugins from the plugin directory
     */
    Uni<List<Plugin>> discoverAndLoadPlugins();

    /**
     * Set the plugin directory
     */
    void setPluginDirectory(String pluginDirectory);

    /**
     * Set the data directory
     */
    void setDataDirectory(String dataDirectory);
}
================================================================================

================================================================================
tech/kayys/silat/plugin/ServiceRegistry.java
Size: 1.1 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin;

import java.util.Optional;

/**
 * Service registry for inter-plugin communication
 * 
 * Plugins can register services that other plugins can discover and use.
 */
public interface ServiceRegistry {
    
    /**
     * Register a service
     * 
     * @param serviceType the service interface type
     * @param service the service implementation
     * @param <T> the service type
     */
    <T> void registerService(Class<T> serviceType, T service);
    
    /**
     * Unregister a service
     * 
     * @param serviceType the service interface type
     * @param <T> the service type
     */
    <T> void unregisterService(Class<T> serviceType);
    
    /**
     * Get a service
     * 
     * @param serviceType the service interface type
     * @param <T> the service type
     * @return the service if registered
     */
    <T> Optional<T> getService(Class<T> serviceType);
    
    /**
     * Check if a service is registered
     * 
     * @param serviceType the service interface type
     * @return true if the service is registered
     */
    boolean hasService(Class<?> serviceType);
}

================================================================================

================================================================================
tech/kayys/silat/plugin/discovery/ServiceDiscoveryPlugin.java
Size: 822 B | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.discovery;

import tech.kayys.silat.plugin.Plugin;
import java.util.Optional;

/**
 * Service Discovery Plugin Interface
 * 
 * Plugins implementing this interface provide dynamic endpoint discovery
 * capabilities for executors.
 * This effectively allows overriding the static endpoint information stored in
 * the registry.
 */
public interface ServiceDiscoveryPlugin extends Plugin {

    /**
     * Discover the endpoint for a given executor ID.
     * 
     * @param executorId The ID of the executor to discover.
     * @return An Optional containing the discovered endpoint (e.g., "host:port" or
     *         "http://host:port"),
     *         or empty if the plugin cannot find an endpoint for this executor.
     */
    Optional<String> discoverEndpoint(String executorId);

}

================================================================================

================================================================================
tech/kayys/silat/plugin/dispatcher/TaskDispatcherPlugin.java
Size: 1.7 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.dispatcher;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.plugin.Plugin;

/**
 * Plugin interface for custom task dispatchers
 * 
 * Task dispatcher plugins can handle dispatching tasks to executors
 * using custom communication protocols beyond the built-in GRPC, Kafka, and REST.
 */
public interface TaskDispatcherPlugin extends Plugin {
    
    /**
     * Check if this dispatcher supports the given executor
     * 
     * @param executor the executor information
     * @return true if this dispatcher can handle the executor
     */
    boolean supports(ExecutorInfo executor);
    
    /**
     * Dispatch a task to an executor
     * 
     * @param task the task to dispatch
     * @param executor the executor to dispatch to
     * @return a Uni that completes when the task is dispatched
     */
    Uni<Void> dispatch(NodeExecutionTask task, ExecutorInfo executor);
    
    /**
     * Get the priority of this dispatcher
     * 
     * Higher priority dispatchers are preferred when multiple dispatchers
     * support the same executor.
     * 
     * @return the priority (default is 0)
     */
    default int getPriority() {
        return 0;
    }
    
    /**
     * Executor information needed for dispatching
     */
    interface ExecutorInfo {
        String executorId();
        String executorType();
        String communicationType();
        String endpoint();
        int timeout();
    }
    
    /**
     * Node execution task information
     */
    interface NodeExecutionTask {
        String runId();
        String nodeId();
        String nodeType();
        java.util.Map<String, Object> inputs();
        int attempt();
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/event/EventListenerPlugin.java
Size: 2.5 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.event;

import tech.kayys.silat.plugin.Plugin;

/**
 * Plugin interface for workflow event listeners
 * 
 * Event listener plugins can react to workflow lifecycle events.
 */
public interface EventListenerPlugin extends Plugin {
    
    /**
     * Called when a workflow run is started
     * 
     * @param event the workflow started event
     */
    default void onWorkflowStarted(WorkflowStartedEvent event) {
        // Default: do nothing
    }
    
    /**
     * Called when a workflow run is completed
     * 
     * @param event the workflow completed event
     */
    default void onWorkflowCompleted(WorkflowCompletedEvent event) {
        // Default: do nothing
    }
    
    /**
     * Called when a workflow run fails
     * 
     * @param event the workflow failed event
     */
    default void onWorkflowFailed(WorkflowFailedEvent event) {
        // Default: do nothing
    }
    
    /**
     * Called when a node is executed
     * 
     * @param event the node executed event
     */
    default void onNodeExecuted(NodeExecutedEvent event) {
        // Default: do nothing
    }
    
    /**
     * Called when a node execution fails
     * 
     * @param event the node failed event
     */
    default void onNodeFailed(NodeFailedEvent event) {
        // Default: do nothing
    }
    
    /**
     * Workflow started event
     */
    interface WorkflowStartedEvent {
        String runId();
        String definitionId();
        java.time.Instant startedAt();
        java.util.Map<String, Object> inputs();
    }
    
    /**
     * Workflow completed event
     */
    interface WorkflowCompletedEvent {
        String runId();
        String definitionId();
        java.time.Instant completedAt();
        java.util.Map<String, Object> outputs();
    }
    
    /**
     * Workflow failed event
     */
    interface WorkflowFailedEvent {
        String runId();
        String definitionId();
        java.time.Instant failedAt();
        String errorMessage();
    }
    
    /**
     * Node executed event
     */
    interface NodeExecutedEvent {
        String runId();
        String nodeId();
        String nodeType();
        java.time.Instant executedAt();
        java.util.Map<String, Object> outputs();
    }
    
    /**
     * Node failed event
     */
    interface NodeFailedEvent {
        String runId();
        String nodeId();
        String nodeType();
        java.time.Instant failedAt();
        String errorMessage();
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/executor/ExecutorPlugin.java
Size: 1.7 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.executor;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.execution.NodeExecutionResult;
import tech.kayys.silat.execution.NodeExecutionTask;
import tech.kayys.silat.plugin.Plugin;

/**
 * Plugin interface for custom task executors
 * 
 * Executor plugins allow dynamic loading of custom task handlers
 * without modifying the core executor runtime.
 * 
 * Example usage:
 * 
 * <pre>
 * {@code
 * public class HttpExecutorPlugin implements ExecutorPlugin {
 *     public String getExecutorType() {
 *         return "http";
 *     }
 * 
 *     public boolean canHandle(NodeExecutionTask task) {
 *         return task.taskType().equals("http-request");
 *     }
 * 
 *     public Uni<NodeExecutionResult> execute(NodeExecutionTask task) {
 *         // Execute HTTP request
 *     }
 * }
 * }
 * </pre>
 */
public interface ExecutorPlugin extends Plugin {

    /**
     * Get the executor type this plugin handles
     * 
     * @return executor type (e.g., "http", "database", "ml-inference")
     */
    String getExecutorType();

    /**
     * Check if this plugin can handle the given task
     * 
     * @param task the task to check
     * @return true if this plugin can handle the task
     */
    boolean canHandle(NodeExecutionTask task);

    /**
     * Execute the task
     * 
     * @param task the task to execute
     * @return execution result wrapped in Uni for reactive execution
     */
    Uni<NodeExecutionResult> execute(NodeExecutionTask task);

    /**
     * Get plugin priority (higher = preferred when multiple plugins can handle a
     * task)
     * 
     * @return priority value, default is 0
     */
    default int getPriority() {
        return 0;
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/impl/DefaultEventBus.java
Size: 2.0 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.impl;

import java.util.List;
import java.util.Map;
import java.util.concurrent.ConcurrentHashMap;
import java.util.concurrent.CopyOnWriteArrayList;
import java.util.function.Consumer;

import jakarta.enterprise.context.ApplicationScoped;
import tech.kayys.silat.plugin.EventBus;
import tech.kayys.silat.plugin.PluginEvent;

/**
 * Default implementation of EventBus
 */
@ApplicationScoped
public class DefaultEventBus implements EventBus {

    private final Map<Class<? extends PluginEvent>, List<Consumer<? extends PluginEvent>>> subscribers = new ConcurrentHashMap<>();

    @Override
    public void publish(PluginEvent event) {
        if (event == null) {
            throw new IllegalArgumentException("Event cannot be null");
        }

        Class<? extends PluginEvent> eventType = event.getClass();
        List<Consumer<? extends PluginEvent>> handlers = subscribers.get(eventType);

        if (handlers != null) {
            for (Consumer<? extends PluginEvent> handler : handlers) {
                try {
                    @SuppressWarnings("unchecked")
                    Consumer<PluginEvent> typedHandler = (Consumer<PluginEvent>) handler;
                    typedHandler.accept(event);
                } catch (Exception e) {
                    // Log but don't fail on handler errors
                    System.err.println("Error in event handler: " + e.getMessage());
                }
            }
        }
    }

    @Override
    public <T extends PluginEvent> Subscription subscribe(Class<T> eventType, Consumer<T> handler) {
        if (eventType == null) {
            throw new IllegalArgumentException("Event type cannot be null");
        }
        if (handler == null) {
            throw new IllegalArgumentException("Handler cannot be null");
        }

        List<Consumer<? extends PluginEvent>> handlers = subscribers.computeIfAbsent(
                eventType,
                k -> new CopyOnWriteArrayList<>());
        handlers.add(handler);

        return () -> handlers.remove(handler);
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/impl/DefaultPluginContext.java
Size: 2.0 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.impl;

import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import org.slf4j.Logger;
import tech.kayys.silat.plugin.PluginContext;
import tech.kayys.silat.plugin.PluginMetadata;
import tech.kayys.silat.plugin.ServiceRegistry;
import tech.kayys.silat.plugin.EventBus;

/**
 * Default implementation of PluginContext
 */
public class DefaultPluginContext implements PluginContext {

    private final PluginMetadata metadata;
    private final Logger logger;
    private final Map<String, String> properties;
    private final ServiceRegistry serviceRegistry;
    private final EventBus eventBus;
    private final String dataDirectory;

    public DefaultPluginContext(
            PluginMetadata metadata,
            Logger logger,
            Map<String, String> properties,
            ServiceRegistry serviceRegistry,
            EventBus eventBus,
            String dataDirectory) {
        this.metadata = metadata;
        this.logger = logger;
        this.properties = new ConcurrentHashMap<>(properties);
        this.serviceRegistry = serviceRegistry;
        this.eventBus = eventBus;
        this.dataDirectory = dataDirectory;
    }

    @Override
    public PluginMetadata getMetadata() {
        return metadata;
    }

    @Override
    public Logger getLogger() {
        return logger;
    }

    @Override
    public Optional<String> getProperty(String key) {
        return Optional.ofNullable(properties.get(key));
    }

    @Override
    public String getProperty(String key, String defaultValue) {
        return properties.getOrDefault(key, defaultValue);
    }

    @Override
    public Map<String, String> getAllProperties() {
        return Map.copyOf(properties);
    }

    @Override
    public ServiceRegistry getServiceRegistry() {
        return serviceRegistry;
    }

    @Override
    public EventBus getEventBus() {
        return eventBus;
    }

    @Override
    public String getDataDirectory() {
        return dataDirectory;
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/impl/DefaultPluginService.java
Size: 3.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.impl;

import java.nio.file.Path;
import java.util.List;
import java.util.Optional;

import io.smallrye.mutiny.Uni;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import tech.kayys.silat.plugin.Plugin;
import tech.kayys.silat.plugin.PluginService;

/**
 * Unified implementation of PluginService that combines all plugin
 * functionality
 */
@ApplicationScoped
@jakarta.enterprise.inject.Typed(PluginService.class)
public class DefaultPluginService implements PluginService {

    @Inject
    PluginManager pluginManager;

    @Inject
    DefaultServiceRegistry serviceRegistry;

    @Inject
    DefaultEventBus eventBus;

    @Override
    public Uni<Plugin> loadPlugin(Path pluginJar) {
        return pluginManager.loadPlugin(pluginJar);
    }

    @Override
    public Uni<Void> registerPlugin(Plugin plugin) {
        return pluginManager.registerPlugin(plugin);
    }

    @Override
    public Uni<Void> startPlugin(String pluginId) {
        return pluginManager.startPlugin(pluginId);
    }

    @Override
    public Uni<Void> stopPlugin(String pluginId) {
        return pluginManager.stopPlugin(pluginId);
    }

    @Override
    public Uni<Void> unloadPlugin(String pluginId) {
        return pluginManager.unloadPlugin(pluginId);
    }

    @Override
    public Uni<Plugin> reloadPlugin(String pluginId, Path pluginJar) {
        return pluginManager.reloadPlugin(pluginId, pluginJar);
    }

    @Override
    public Optional<Plugin> getPlugin(String pluginId) {
        return pluginManager.getPlugin(pluginId);
    }

    @Override
    public List<Plugin> getAllPlugins() {
        return pluginManager.getAllPlugins();
    }

    @Override
    public <T extends Plugin> List<T> getPluginsByType(Class<T> pluginType) {
        return pluginManager.getPluginsByType(pluginType);
    }

    @Override
    public Uni<List<Plugin>> discoverAndLoadPlugins() {
        return pluginManager.discoverAndLoadPlugins();
    }

    @Override
    public void setPluginDirectory(String pluginDirectory) {
        pluginManager.setPluginDirectory(pluginDirectory);
    }

    @Override
    public void setDataDirectory(String dataDirectory) {
        pluginManager.setDataDirectory(dataDirectory);
    }

    @Override
    public <T> void registerService(Class<T> serviceType, T service) {
        serviceRegistry.registerService(serviceType, service);
    }

    @Override
    public <T> void unregisterService(Class<T> serviceType) {
        serviceRegistry.unregisterService(serviceType);
    }

    @Override
    public <T> Optional<T> getService(Class<T> serviceType) {
        return serviceRegistry.getService(serviceType);
    }

    @Override
    public boolean hasService(Class<?> serviceType) {
        return serviceRegistry.hasService(serviceType);
    }

    @Override
    public void publish(tech.kayys.silat.plugin.PluginEvent event) {
        eventBus.publish(event);
    }

    @Override
    public <T extends tech.kayys.silat.plugin.PluginEvent> tech.kayys.silat.plugin.EventBus.Subscription subscribe(
            Class<T> eventType, java.util.function.Consumer<T> handler) {
        return eventBus.subscribe(eventType, handler);
    }
}
================================================================================

================================================================================
tech/kayys/silat/plugin/impl/DefaultServiceRegistry.java
Size: 1.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.impl;

import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import jakarta.enterprise.context.ApplicationScoped;
import tech.kayys.silat.plugin.ServiceRegistry;

/**
 * Default implementation of ServiceRegistry
 */
@ApplicationScoped
public class DefaultServiceRegistry implements ServiceRegistry {

    private final Map<Class<?>, Object> services = new ConcurrentHashMap<>();

    @Override
    public <T> void registerService(Class<T> serviceType, T service) {
        if (serviceType == null) {
            throw new IllegalArgumentException("Service type cannot be null");
        }
        if (service == null) {
            throw new IllegalArgumentException("Service cannot be null");
        }
        services.put(serviceType, service);
    }

    @Override
    public <T> void unregisterService(Class<T> serviceType) {
        services.remove(serviceType);
    }

    @Override
    @SuppressWarnings("unchecked")
    public <T> Optional<T> getService(Class<T> serviceType) {
        return Optional.ofNullable((T) services.get(serviceType));
    }

    @Override
    public boolean hasService(Class<?> serviceType) {
        return services.containsKey(serviceType);
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/impl/PluginClassLoader.java
Size: 2.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.impl;

import java.io.IOException;
import java.net.URL;
import java.net.URLClassLoader;
import java.nio.file.Path;
import java.util.ArrayList;
import java.util.List;

/**
 * Custom classloader for plugin isolation
 *
 * Uses parent-last delegation for plugin classes to ensure plugins
 * can use their own versions of dependencies.
 */
public class PluginClassLoader extends URLClassLoader {

    private final List<String> sharedPackages;

    public PluginClassLoader(Path pluginJar, ClassLoader parent) {
        super(new URL[] { toURL(pluginJar) }, parent);
        this.sharedPackages = new ArrayList<>();
        // Always share plugin API classes
        this.sharedPackages.add("tech.kayys.silat.plugin");
    }

    /**
     * Add a package to be shared with the parent classloader
     */
    public void addSharedPackage(String packageName) {
        sharedPackages.add(packageName);
    }

    @Override
    protected Class<?> loadClass(String name, boolean resolve) throws ClassNotFoundException {
        synchronized (getClassLoadingLock(name)) {
            // Check if already loaded
            Class<?> c = findLoadedClass(name);
            if (c != null) {
                return c;
            }

            // Check if this is a shared package
            if (isSharedPackage(name)) {
                return super.loadClass(name, resolve);
            }

            // Try to load from plugin first (parent-last)
            try {
                c = findClass(name);
                if (resolve) {
                    resolveClass(c);
                }
                return c;
            } catch (ClassNotFoundException e) {
                // Fall back to parent
                return super.loadClass(name, resolve);
            }
        }
    }

    private boolean isSharedPackage(String className) {
        for (String pkg : sharedPackages) {
            if (className.startsWith(pkg)) {
                return true;
            }
        }
        return false;
    }

    private static URL toURL(Path path) {
        try {
            return path.toUri().toURL();
        } catch (IOException e) {
            throw new RuntimeException("Failed to convert path to URL: " + path, e);
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/impl/PluginManager.java
Size: 12.8 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.impl;

import java.io.IOException;
import java.nio.file.DirectoryStream;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Optional;
import java.util.ServiceLoader;
import java.util.concurrent.ConcurrentHashMap;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import io.smallrye.mutiny.Uni;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;
import tech.kayys.silat.plugin.Plugin;
import tech.kayys.silat.plugin.PluginContext;
import tech.kayys.silat.plugin.PluginException;
import tech.kayys.silat.plugin.PluginMetadata;

/**
 * Central plugin manager for loading, managing, and unloading plugins
 */
@ApplicationScoped
public class PluginManager {

    private static final Logger LOG = LoggerFactory.getLogger(PluginManager.class);

    private final tech.kayys.silat.plugin.impl.PluginRegistry registry = new PluginRegistry();
    private final Map<String, PluginClassLoader> classLoaders = new ConcurrentHashMap<>();

    @Inject
    tech.kayys.silat.plugin.ServiceRegistry serviceRegistry;

    @Inject
    tech.kayys.silat.plugin.EventBus eventBus;

    private String pluginDirectory = "/opt/silat/plugins";
    private String dataDirectory = "/opt/silat/plugin-data";

    /**
     * Load a plugin from a JAR file
     */
    public Uni<Plugin> loadPlugin(Path pluginJar) {
        return Uni.createFrom().item(() -> {
            try {
                LOG.info("Loading plugin from: {}", pluginJar);

                // Create plugin classloader
                PluginClassLoader classLoader = new PluginClassLoader(pluginJar, getClass().getClassLoader());

                // Use ServiceLoader to discover plugin
                ServiceLoader<Plugin> loader = ServiceLoader.load(Plugin.class, classLoader);
                Optional<Plugin> pluginOpt = loader.findFirst();

                if (pluginOpt.isEmpty()) {
                    throw new RuntimeException("No plugin found in JAR: " + pluginJar);
                }

                Plugin plugin = pluginOpt.get();
                PluginMetadata metadata = plugin.getMetadata();

                // Check if already loaded
                if (registry.isRegistered(metadata.id())) {
                    throw new RuntimeException("Plugin already loaded: " + metadata.id());
                }

                // Create plugin context
                String pluginDataDir = dataDirectory + "/" + metadata.id();
                createDirectoryIfNotExists(Paths.get(pluginDataDir));

                PluginContext context = new DefaultPluginContext(
                        metadata,
                        LoggerFactory.getLogger("plugin." + metadata.id()),
                        metadata.properties(),
                        serviceRegistry,
                        eventBus,
                        pluginDataDir);

                // Initialize plugin
                plugin.initialize(context);

                // Register plugin
                tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin loadedPlugin = new tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin(
                        plugin, metadata, classLoader);
                loadedPlugin.setState(tech.kayys.silat.plugin.impl.PluginRegistry.PluginState.INITIALIZED);
                registry.register(loadedPlugin);
                classLoaders.put(metadata.id(), classLoader);

                LOG.info("Plugin loaded successfully: {} v{}", metadata.name(), metadata.version());
                return plugin;

            } catch (PluginException e) {
                LOG.error("Failed to initialize plugin", e);
                throw new RuntimeException("Failed to initialize plugin: " + e.getMessage(), e);
            } catch (Exception e) {
                LOG.error("Failed to load plugin from: {}", pluginJar, e);
                throw new RuntimeException("Failed to load plugin: " + e.getMessage(), e);
            }
        });
    }

    /**
     * Register a plugin instance directly (programmatic registration)
     */
    public Uni<Void> registerPlugin(Plugin plugin) {
        return Uni.createFrom().item(() -> {
            try {
                PluginMetadata metadata = plugin.getMetadata();
                LOG.info("Registering plugin: {} v{}", metadata.name(), metadata.version());

                if (registry.isRegistered(metadata.id())) {
                    throw new RuntimeException("Plugin already registered: " + metadata.id());
                }

                // Create plugin context
                String pluginDataDir = dataDirectory + "/" + metadata.id();
                createDirectoryIfNotExists(Paths.get(pluginDataDir));

                PluginContext context = new DefaultPluginContext(
                        metadata,
                        LoggerFactory.getLogger("plugin." + metadata.id()),
                        metadata.properties(),
                        serviceRegistry,
                        eventBus,
                        pluginDataDir);

                // Initialize plugin
                plugin.initialize(context);

                // Register plugin
                tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin loadedPlugin = new tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin(
                        plugin, metadata, null); // No dedicated classloader for programmatic plugins
                loadedPlugin.setState(tech.kayys.silat.plugin.impl.PluginRegistry.PluginState.INITIALIZED);
                registry.register(loadedPlugin);

                return null;
            } catch (PluginException e) {
                LOG.error("Failed to initialize registered plugin", e);
                throw new RuntimeException("Failed to initialize registered plugin: " + e.getMessage(), e);
            } catch (Exception e) {
                LOG.error("Failed to register plugin", e);
                throw new RuntimeException("Failed to register plugin: " + e.getMessage(), e);
            }
        });
    }

    /**
     * Start a plugin
     */
    public Uni<Void> startPlugin(String pluginId) {
        return Uni.createFrom().item(() -> {
            try {
                Optional<tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin> loadedOpt = registry
                        .getPlugin(pluginId);
                if (loadedOpt.isEmpty()) {
                    throw new RuntimeException("Plugin not found: " + pluginId);
                }

                tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin loaded = loadedOpt.get();
                loaded.getPlugin().start();
                loaded.setState(tech.kayys.silat.plugin.impl.PluginRegistry.PluginState.STARTED);
                LOG.info("Plugin started: {}", pluginId);
                return null;
            } catch (PluginException e) {
                LOG.error("Failed to start plugin: {}", pluginId, e);
                throw new RuntimeException("Failed to start plugin: " + e.getMessage(), e);
            }
        });
    }

    /**
     * Stop a plugin
     */
    public Uni<Void> stopPlugin(String pluginId) {
        return Uni.createFrom().item(() -> {
            try {
                Optional<tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin> loadedOpt = registry
                        .getPlugin(pluginId);
                if (loadedOpt.isEmpty()) {
                    throw new RuntimeException("Plugin not found: " + pluginId);
                }

                tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin loaded = loadedOpt.get();
                loaded.getPlugin().stop();
                loaded.setState(tech.kayys.silat.plugin.impl.PluginRegistry.PluginState.STOPPED);
                LOG.info("Plugin stopped: {}", pluginId);
                return null;
            } catch (PluginException e) {
                LOG.error("Failed to stop plugin: {}", pluginId, e);
                throw new RuntimeException("Failed to stop plugin: " + e.getMessage(), e);
            }
        });
    }

    /**
     * Unload a plugin
     */
    public Uni<Void> unloadPlugin(String pluginId) {
        return stopPlugin(pluginId)
                .onFailure().recoverWithNull()
                .chain(() -> Uni.createFrom().item(() -> {
                    registry.unregister(pluginId);
                    PluginClassLoader classLoader = classLoaders.remove(pluginId);
                    if (classLoader != null) {
                        try {
                            classLoader.close();
                        } catch (IOException e) {
                            LOG.warn("Failed to close classloader for plugin: {}", pluginId, e);
                        }
                    }
                    LOG.info("Plugin unloaded: {}", pluginId);
                    return null;
                }));
    }

    /**
     * Reload a plugin (hot-reload)
     */
    public Uni<Plugin> reloadPlugin(String pluginId, Path pluginJar) {
        return unloadPlugin(pluginId)
                .chain(() -> loadPlugin(pluginJar))
                .chain(plugin -> startPlugin(pluginId).replaceWith(plugin));
    }

    /**
     * Get a plugin by ID
     */
    public Optional<Plugin> getPlugin(String pluginId) {
        return registry.getPlugin(pluginId).map(tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin::getPlugin);
    }

    /**
     * Get all loaded plugins
     */
    public List<Plugin> getAllPlugins() {
        return registry.getAllPlugins().values().stream()
                .map(tech.kayys.silat.plugin.impl.PluginRegistry.LoadedPlugin::getPlugin)
                .toList();
    }

    /**
     * Get plugins by type
     */
    @SuppressWarnings("unchecked")
    public <T extends Plugin> List<T> getPluginsByType(Class<T> pluginType) {
        return registry.getAllPlugins().values().stream()
                .map(PluginRegistry.LoadedPlugin::getPlugin)
                .filter(pluginType::isInstance)
                .map(p -> (T) p)
                .toList();
    }

    /**
     * Discover and load all plugins from the plugin directory and classpath
     */
    public Uni<List<Plugin>> discoverAndLoadPlugins() {
        return Uni.createFrom().item(() -> {
            List<Plugin> loadedPlugins = new ArrayList<>();

            // 1. Load from classpath using ServiceLoader
            LOG.info("Discovering plugins from classpath...");
            ServiceLoader<Plugin> loader = ServiceLoader.load(Plugin.class);
            for (Plugin plugin : loader) {
                try {
                    if (!registry.isRegistered(plugin.getMetadata().id())) {
                        LOG.info("Discovered classpath plugin: {}", plugin.getMetadata().id());
                        registerPlugin(plugin).await().indefinitely();
                        startPlugin(plugin.getMetadata().id()).await().indefinitely();
                        loadedPlugins.add(plugin);
                    }
                } catch (Exception e) {
                    LOG.error("Failed to load classpath plugin: {}", plugin.getClass().getName(), e);
                }
            }

            // 2. Load from plugin directory
            Path pluginDir = Paths.get(pluginDirectory);
            if (Files.exists(pluginDir)) {
                LOG.info("Scanning plugin directory: {}", pluginDirectory);
                try (DirectoryStream<Path> stream = Files.newDirectoryStream(pluginDir, "*.jar")) {
                    for (Path jarFile : stream) {
                        try {
                            Plugin plugin = loadPlugin(jarFile).await().indefinitely();
                            startPlugin(plugin.getMetadata().id()).await().indefinitely();
                            loadedPlugins.add(plugin);
                        } catch (Exception e) {
                            LOG.error("Failed to load plugin from: {}", jarFile, e);
                        }
                    }
                } catch (IOException e) {
                    LOG.error("Failed to scan plugin directory", e);
                }
            }

            LOG.info("Total plugins discovered and loaded: {}", loadedPlugins.size());
            return loadedPlugins;
        });
    }

    /**
     * Set the plugin directory
     */
    public void setPluginDirectory(String pluginDirectory) {
        this.pluginDirectory = pluginDirectory;
    }

    /**
     * Set the data directory
     */
    public void setDataDirectory(String dataDirectory) {
        this.dataDirectory = dataDirectory;
    }

    /**
     * Get the plugin registry (for internal use)
     */
    public PluginRegistry getRegistry() {
        return registry;
    }

    private void createDirectoryIfNotExists(Path dir) {
        try {
            if (!Files.exists(dir)) {
                Files.createDirectories(dir);
            }
        } catch (IOException e) {
            LOG.warn("Failed to create directory: {}", dir, e);
        }
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/impl/PluginRegistry.java
Size: 3.0 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.impl;

import java.util.Map;
import java.util.Optional;
import java.util.concurrent.ConcurrentHashMap;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import tech.kayys.silat.plugin.Plugin;
import tech.kayys.silat.plugin.PluginMetadata;

/**
 * Registry for managing loaded plugins
 */
public class PluginRegistry {

    private static final Logger LOG = LoggerFactory.getLogger(PluginRegistry.class);

    private final Map<String, LoadedPlugin> plugins = new ConcurrentHashMap<>();

    /**
     * Register a loaded plugin
     */
    public void register(LoadedPlugin plugin) {
        String pluginId = plugin.getMetadata().id();
        if (plugins.containsKey(pluginId)) {
            throw new IllegalStateException("Plugin already registered: " + pluginId);
        }
        plugins.put(pluginId, plugin);
        LOG.info("Registered plugin: {} v{}", plugin.getMetadata().name(), plugin.getMetadata().version());
    }

    /**
     * Unregister a plugin
     */
    public void unregister(String pluginId) {
        LoadedPlugin plugin = plugins.remove(pluginId);
        if (plugin != null) {
            LOG.info("Unregistered plugin: {}", pluginId);
        }
    }

    /**
     * Get a plugin by ID
     */
    public Optional<LoadedPlugin> getPlugin(String pluginId) {
        return Optional.ofNullable(plugins.get(pluginId));
    }

    /**
     * Get all loaded plugins
     */
    public Map<String, LoadedPlugin> getAllPlugins() {
        return Map.copyOf(plugins);
    }

    /**
     * Check if a plugin is registered
     */
    public boolean isRegistered(String pluginId) {
        return plugins.containsKey(pluginId);
    }

    /**
     * Get the number of loaded plugins
     */
    public int getPluginCount() {
        return plugins.size();
    }

    /**
     * Loaded plugin information
     */
    public static class LoadedPlugin {
        private final Plugin plugin;
        private final PluginMetadata metadata;
        private final tech.kayys.silat.plugin.impl.PluginClassLoader classLoader;
        private PluginState state;

        public LoadedPlugin(Plugin plugin, PluginMetadata metadata, tech.kayys.silat.plugin.impl.PluginClassLoader classLoader) {
            this.plugin = plugin;
            this.metadata = metadata;
            this.classLoader = classLoader;
            this.state = PluginState.LOADED;
        }

        public Plugin getPlugin() {
            return plugin;
        }

        public PluginMetadata getMetadata() {
            return metadata;
        }

        public tech.kayys.silat.plugin.impl.PluginClassLoader getClassLoader() {
            return classLoader;
        }

        public PluginState getState() {
            return state;
        }

        public void setState(PluginState state) {
            this.state = state;
        }
    }

    /**
     * Plugin lifecycle state
     */
    public enum PluginState {
        LOADED,
        INITIALIZED,
        STARTED,
        STOPPED,
        FAILED
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/interceptor/ExecutionInterceptorPlugin.java
Size: 2.0 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.interceptor;

import io.smallrye.mutiny.Uni;
import tech.kayys.silat.plugin.Plugin;

/**
 * Plugin interface for execution interceptors
 * 
 * Execution interceptor plugins can hook into the task execution lifecycle
 * to perform actions before, after, or on error during task execution.
 */
public interface ExecutionInterceptorPlugin extends Plugin {
    
    /**
     * Called before a task is executed
     * 
     * @param task the task about to be executed
     * @return a Uni that completes when pre-processing is done
     */
    default Uni<Void> beforeExecution(TaskContext task) {
        return Uni.createFrom().voidItem();
    }
    
    /**
     * Called after a task is successfully executed
     * 
     * @param task the task that was executed
     * @param result the execution result
     * @return a Uni that completes when post-processing is done
     */
    default Uni<Void> afterExecution(TaskContext task, ExecutionResult result) {
        return Uni.createFrom().voidItem();
    }
    
    /**
     * Called when a task execution fails
     * 
     * @param task the task that failed
     * @param error the error that occurred
     * @return a Uni that completes when error handling is done
     */
    default Uni<Void> onError(TaskContext task, Throwable error) {
        return Uni.createFrom().voidItem();
    }
    
    /**
     * Get the order of this interceptor
     * 
     * Lower order interceptors are executed first.
     * 
     * @return the order (default is 0)
     */
    default int getOrder() {
        return 0;
    }
    
    /**
     * Task context information
     */
    interface TaskContext {
        String runId();
        String nodeId();
        String nodeType();
        java.util.Map<String, Object> inputs();
        int attempt();
    }
    
    /**
     * Execution result information
     */
    interface ExecutionResult {
        boolean isSuccess();
        java.util.Map<String, Object> outputs();
        String errorMessage();
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/transformer/DataTransformerPlugin.java
Size: 1.2 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.transformer;

import tech.kayys.silat.plugin.Plugin;
import java.util.Map;

/**
 * Plugin interface for data transformers
 * 
 * Transformer plugins can transform input/output data for workflow nodes.
 */
public interface DataTransformerPlugin extends Plugin {
    
    /**
     * Check if this transformer supports the given node type
     * 
     * @param nodeType the node type
     * @return true if this transformer can handle the node type
     */
    boolean supports(String nodeType);
    
    /**
     * Transform input data before task execution
     * 
     * @param input the input data
     * @param node the node definition
     * @return the transformed input data
     */
    Map<String, Object> transformInput(Map<String, Object> input, NodeContext node);
    
    /**
     * Transform output data after task execution
     * 
     * @param output the output data
     * @param node the node definition
     * @return the transformed output data
     */
    Map<String, Object> transformOutput(Map<String, Object> output, NodeContext node);
    
    /**
     * Node context information
     */
    interface NodeContext {
        String nodeId();
        String nodeType();
        Map<String, Object> configuration();
    }
}

================================================================================

================================================================================
tech/kayys/silat/plugin/validator/WorkflowValidatorPlugin.java
Size: 1.6 KB | Modified: 2026-01-17 05:59:41
--------------------------------------------------------------------------------
package tech.kayys.silat.plugin.validator;

import tech.kayys.silat.plugin.Plugin;
import java.util.List;

/**
 * Plugin interface for workflow validators
 * 
 * Validator plugins can add custom validation rules for workflow definitions.
 */
public interface WorkflowValidatorPlugin extends Plugin {
    
    /**
     * Validate a workflow definition
     * 
     * @param definition the workflow definition to validate
     * @return list of validation errors (empty if valid)
     */
    List<ValidationError> validate(WorkflowDefinition definition);
    
    /**
     * Get the validation rules provided by this plugin
     * 
     * @return list of validation rule descriptions
     */
    List<String> getValidationRules();
    
    /**
     * Workflow definition information
     */
    interface WorkflowDefinition {
        String definitionId();
        String name();
        String version();
        List<NodeDefinition> nodes();
        List<Transition> transitions();
    }
    
    /**
     * Node definition information
     */
    interface NodeDefinition {
        String nodeId();
        String nodeType();
        java.util.Map<String, Object> configuration();
    }
    
    /**
     * Transition information
     */
    interface Transition {
        String fromNodeId();
        String toNodeId();
        String condition();
    }
    
    /**
     * Validation error
     */
    record ValidationError(
        String rule,
        String message,
        String location,
        Severity severity
    ) {
        public enum Severity {
            ERROR, WARNING, INFO
        }
    }
}

================================================================================


