Coto Output
Generated: 2026-01-17 09:27:22
Files: 12 | Directories: 10 | Total Size: 48.2 KB


================================================================================
tech/kayys/silat/runtime/DbInitializer.java
Size: 1.8 KB | Modified: 2026-01-15 12:52:37
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
tech/kayys/silat/runtime/resource/CallbackResource.java
Size: 964 B | Modified: 2026-01-15 12:52:37
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
Size: 4.7 KB | Modified: 2026-01-15 12:52:37
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
Size: 2.2 KB | Modified: 2026-01-15 12:52:37
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
Size: 2.2 KB | Modified: 2026-01-15 12:52:37
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
Size: 4.1 KB | Modified: 2026-01-15 12:52:37
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
tech/kayys/silat/runtime/standalone/SilatStandaloneRuntime.java
Size: 665 B | Modified: 2026-01-16 06:22:41
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.standalone;

import io.quarkus.runtime.Quarkus;
import io.quarkus.runtime.annotations.QuarkusMain;

/**
 * Main entry point for the Silat Standalone Runtime
 * This is a pure server runtime that hosts the workflow engine
 */
@QuarkusMain
public class SilatStandaloneRuntime {

    public static void main(String[] args) {
        System.out.println("Starting Silat Standalone Runtime Server...");
        System.out.println("Server will listen on HTTP port (configured in application.properties)");
        System.out.println("gRPC server will listen on port (configured in application.properties)");

        Quarkus.run();
    }
}
================================================================================

================================================================================
tech/kayys/silat/runtime/standalone/plugin/PluginConfigurationService.java
Size: 3.8 KB | Modified: 2026-01-15 18:01:34
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.standalone.plugin;

import io.quarkus.logging.Log;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.inject.Inject;

import java.io.*;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.Properties;

/**
 * Manages plugin configuration settings
 */
@ApplicationScoped
public class PluginConfigurationService {

    @Inject
    PluginManager pluginManager;

    private static final String PLUGIN_CONFIG_FILE = "plugin-config.properties";

    /**
     * Gets the configuration file path for a specific plugin
     */
    public Path getPluginConfigPath(String pluginName) {
        String pluginsDir = pluginManager.getPluginsDirectory();
        return Paths.get(pluginsDir, pluginName, PLUGIN_CONFIG_FILE);
    }

    /**
     * Loads configuration for a specific plugin
     */
    public Properties loadPluginConfig(String pluginName) {
        Properties props = new Properties();
        Path configPath = getPluginConfigPath(pluginName);

        if (Files.exists(configPath)) {
            try (InputStream input = Files.newInputStream(configPath)) {
                props.load(input);
                Log.infof("Loaded configuration for plugin: %s", pluginName);
            } catch (IOException e) {
                Log.errorf("Failed to load configuration for plugin %s: %s", pluginName, e.getMessage());
            }
        }

        return props;
    }

    /**
     * Saves configuration for a specific plugin
     */
    public boolean savePluginConfig(String pluginName, Properties config) {
        Path configPath = getPluginConfigPath(pluginName);
        Path pluginDir = configPath.getParent();

        try {
            // Create plugin directory if it doesn't exist
            if (pluginDir != null && !Files.exists(pluginDir)) {
                Files.createDirectories(pluginDir);
            }

            try (OutputStream output = Files.newOutputStream(configPath)) {
                config.store(output, "Plugin configuration for " + pluginName);
                Log.infof("Saved configuration for plugin: %s", pluginName);
                return true;
            }
        } catch (IOException e) {
            Log.errorf("Failed to save configuration for plugin %s: %s", pluginName, e.getMessage());
            return false;
        }
    }

    /**
     * Updates a specific configuration property for a plugin
     */
    public boolean updatePluginConfigProperty(String pluginName, String key, String value) {
        Properties config = loadPluginConfig(pluginName);
        config.setProperty(key, value);
        return savePluginConfig(pluginName, config);
    }

    /**
     * Gets a specific configuration property for a plugin
     */
    public String getPluginConfigProperty(String pluginName, String key, String defaultValue) {
        Properties config = loadPluginConfig(pluginName);
        return config.getProperty(key, defaultValue);
    }

    /**
     * Removes a configuration property for a plugin
     */
    public boolean removePluginConfigProperty(String pluginName, String key) {
        Properties config = loadPluginConfig(pluginName);
        if (config.containsKey(key)) {
            config.remove(key);
            return savePluginConfig(pluginName, config);
        }
        return true; // Property didn't exist anyway
    }

    /**
     * Creates a default configuration template for a plugin
     */
    public boolean createDefaultConfigTemplate(String pluginName) {
        Properties defaultProps = new Properties();
        defaultProps.setProperty("enabled", "true");
        defaultProps.setProperty("auto-start", "true");
        defaultProps.setProperty("thread-pool-size", "5");
        defaultProps.setProperty("timeout-seconds", "30");

        return savePluginConfig(pluginName, defaultProps);
    }
}
================================================================================

================================================================================
tech/kayys/silat/runtime/standalone/plugin/PluginInfo.java
Size: 1.4 KB | Modified: 2026-01-15 17:57:53
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.standalone.plugin;

import java.net.URLClassLoader;

/**
 * Represents information about a loaded plugin
 */
public class PluginInfo {
    private final String name;
    private final String version;
    private final String fileName;
    private final String filePath;
    private final Class<?> pluginClass;
    private final URLClassLoader classLoader;
    private boolean enabled;

    public PluginInfo(String name, String version, String fileName, String filePath, 
                      Class<?> pluginClass, URLClassLoader classLoader, boolean enabled) {
        this.name = name;
        this.version = version;
        this.fileName = fileName;
        this.filePath = filePath;
        this.pluginClass = pluginClass;
        this.classLoader = classLoader;
        this.enabled = enabled;
    }

    public String getName() {
        return name;
    }

    public String getVersion() {
        return version;
    }

    public String getFileName() {
        return fileName;
    }

    public String getFilePath() {
        return filePath;
    }

    public Class<?> getPluginClass() {
        return pluginClass;
    }

    public URLClassLoader getClassLoader() {
        return classLoader;
    }

    public boolean isEnabled() {
        return enabled;
    }

    public void setEnabled(boolean enabled) {
        this.enabled = enabled;
    }
}
================================================================================

================================================================================
tech/kayys/silat/runtime/standalone/plugin/PluginManager.java
Size: 7.3 KB | Modified: 2026-01-15 17:57:38
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.standalone.plugin;

import io.quarkus.runtime.StartupEvent;
import jakarta.enterprise.context.ApplicationScoped;
import jakarta.enterprise.event.Observes;
import jakarta.inject.Inject;
import org.eclipse.microprofile.config.inject.ConfigProperty;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

import java.io.File;
import java.io.IOException;
import java.net.URL;
import java.net.URLClassLoader;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.*;
import java.util.jar.JarEntry;
import java.util.jar.JarFile;

/**
 * Manages plugin loading, unloading, and lifecycle
 */
@ApplicationScoped
public class PluginManager {

    private static final Logger LOG = LoggerFactory.getLogger(PluginManager.class);

    @Inject
    @ConfigProperty(name = "silat.plugins.directory", defaultValue = "./plugins")
    String pluginsDirectory;

    @Inject
    @ConfigProperty(name = "silat.plugins.auto-discover", defaultValue = "true")
    boolean autoDiscoverPlugins;

    private final Map<String, PluginInfo> loadedPlugins = new HashMap<>();
    private final List<ClassLoader> pluginClassLoaders = new ArrayList<>();

    void onStart(@Observes StartupEvent ev) {
        LOG.info("Initializing Plugin Manager with directory: {}", pluginsDirectory);
        
        // Create plugins directory if it doesn't exist
        Path pluginDir = Paths.get(pluginsDirectory);
        if (!Files.exists(pluginDir)) {
            try {
                Files.createDirectories(pluginDir);
                LOG.info("Created plugins directory: {}", pluginDir.toAbsolutePath());
            } catch (IOException e) {
                LOG.error("Failed to create plugins directory: {}", e.getMessage());
            }
        }

        if (autoDiscoverPlugins) {
            scanAndLoadPlugins();
        }
    }

    /**
     * Scans the plugins directory and loads all available plugins
     */
    public void scanAndLoadPlugins() {
        LOG.info("Scanning for plugins in directory: {}", pluginsDirectory);
        
        try {
            Files.walk(Paths.get(pluginsDirectory))
                    .filter(path -> path.toString().endsWith(".jar"))
                    .forEach(this::loadPlugin);
        } catch (IOException e) {
            LOG.error("Error scanning plugins directory: {}", e.getMessage());
        }
    }

    /**
     * Loads a plugin from the specified JAR file
     */
    public synchronized boolean loadPlugin(Path jarPath) {
        String fileName = jarPath.getFileName().toString();
        LOG.info("Loading plugin: {}", fileName);

        try (JarFile jarFile = new JarFile(jarPath.toFile())) {
            // Check if plugin is already loaded
            if (loadedPlugins.containsKey(fileName)) {
                LOG.warn("Plugin {} is already loaded", fileName);
                return false;
            }

            // Extract plugin metadata from manifest
            String pluginName = jarFile.getManifest().getMainAttributes().getValue("Plugin-Name");
            String pluginVersion = jarFile.getManifest().getMainAttributes().getValue("Plugin-Version");
            String pluginClass = jarFile.getManifest().getMainAttributes().getValue("Plugin-Class");

            if (pluginName == null || pluginClass == null) {
                LOG.error("Plugin {} is missing required manifest attributes", fileName);
                return false;
            }

            // Create class loader for the plugin
            URL jarUrl = jarPath.toUri().toURL();
            URLClassLoader classLoader = new URLClassLoader(new URL[]{jarUrl}, 
                    Thread.currentThread().getContextClassLoader());
            
            // Load the plugin class
            Class<?> pluginClazz = classLoader.loadClass(pluginClass);
            
            // Store plugin info
            PluginInfo pluginInfo = new PluginInfo(
                    pluginName,
                    pluginVersion != null ? pluginVersion : "unknown",
                    fileName,
                    jarPath.toAbsolutePath().toString(),
                    pluginClazz,
                    classLoader,
                    true // enabled by default
            );

            loadedPlugins.put(fileName, pluginInfo);
            pluginClassLoaders.add(classLoader);

            LOG.info("Successfully loaded plugin: {} ({})", pluginName, fileName);
            return true;

        } catch (Exception e) {
            LOG.error("Failed to load plugin {}: {}", fileName, e.getMessage());
            return false;
        }
    }

    /**
     * Unloads a plugin by name
     */
    public synchronized boolean unloadPlugin(String pluginFileName) {
        PluginInfo pluginInfo = loadedPlugins.get(pluginFileName);
        if (pluginInfo == null) {
            LOG.warn("Plugin {} is not loaded", pluginFileName);
            return false;
        }

        try {
            // Remove from loaded plugins
            loadedPlugins.remove(pluginFileName);
            
            // Remove class loader
            pluginClassLoaders.remove(pluginInfo.getClassLoader());
            
            // Attempt to close the class loader (Java 9+ feature)
            if (pluginInfo.getClassLoader() instanceof URLClassLoader) {
                try {
                    ((URLClassLoader) pluginInfo.getClassLoader()).close();
                } catch (IOException e) {
                    LOG.warn("Could not close class loader for plugin {}: {}", 
                            pluginFileName, e.getMessage());
                }
            }

            LOG.info("Successfully unloaded plugin: {}", pluginInfo.getName());
            return true;

        } catch (Exception e) {
            LOG.error("Failed to unload plugin {}: {}", pluginFileName, e.getMessage());
            return false;
        }
    }

    /**
     * Enables a plugin
     */
    public boolean enablePlugin(String pluginFileName) {
        PluginInfo pluginInfo = loadedPlugins.get(pluginFileName);
        if (pluginInfo == null) {
            LOG.warn("Cannot enable plugin {}: plugin not loaded", pluginFileName);
            return false;
        }

        pluginInfo.setEnabled(true);
        LOG.info("Enabled plugin: {}", pluginInfo.getName());
        return true;
    }

    /**
     * Disables a plugin
     */
    public boolean disablePlugin(String pluginFileName) {
        PluginInfo pluginInfo = loadedPlugins.get(pluginFileName);
        if (pluginInfo == null) {
            LOG.warn("Cannot disable plugin {}: plugin not loaded", pluginFileName);
            return false;
        }

        pluginInfo.setEnabled(false);
        LOG.info("Disabled plugin: {}", pluginInfo.getName());
        return true;
    }

    /**
     * Gets information about all loaded plugins
     */
    public List<PluginInfo> getAllPlugins() {
        return new ArrayList<>(loadedPlugins.values());
    }

    /**
     * Gets information about a specific plugin
     */
    public PluginInfo getPlugin(String pluginFileName) {
        return loadedPlugins.get(pluginFileName);
    }

    /**
     * Gets the plugins directory path
     */
    public String getPluginsDirectory() {
        return pluginsDirectory;
    }

    /**
     * Refreshes the plugin list by rescanning the directory
     */
    public void refreshPlugins() {
        LOG.info("Refreshing plugins...");
        scanAndLoadPlugins();
    }
}
================================================================================

================================================================================
tech/kayys/silat/runtime/standalone/plugin/PluginResource.java
Size: 17.2 KB | Modified: 2026-01-16 07:54:05
--------------------------------------------------------------------------------
package tech.kayys.silat.runtime.standalone.plugin;

import io.quarkus.logging.Log;
import jakarta.enterprise.context.RequestScoped;
import jakarta.inject.Inject;
import jakarta.ws.rs.*;
import jakarta.ws.rs.core.MediaType;
import jakarta.ws.rs.core.Response;
import org.jboss.resteasy.annotations.providers.multipart.PartType;
import org.jboss.resteasy.annotations.providers.multipart.MultipartForm;
import org.jboss.resteasy.plugins.providers.multipart.InputPart;
import org.jboss.resteasy.plugins.providers.multipart.MultipartFormDataInput;

import java.io.File;
import java.io.FileOutputStream;
import java.io.InputStream;
import java.nio.file.Files;
import java.nio.file.Paths;
import java.util.List;
import java.util.Map;
import java.util.Properties;

/**
 * REST endpoint for plugin management and upload
 */
@jakarta.ws.rs.Path("/api/plugins")
@Produces(MediaType.APPLICATION_JSON)
@Consumes(MediaType.MULTIPART_FORM_DATA)
@RequestScoped
public class PluginResource {

    @Inject
    PluginManager pluginManager;

    @Inject
    PluginConfigurationService pluginConfigService;

    @POST
    @jakarta.ws.rs.Path("/upload")
    @Consumes(MediaType.MULTIPART_FORM_DATA)
    public Response uploadPlugin(MultipartFormDataInput input) {
        try {
            // Extract uploaded file and filename from multipart input
            Map<String, List<InputPart>> uploadForm = input.getFormDataMap();

            List<InputPart> fileParts = uploadForm.get("uploadedInputStream");
            if (fileParts == null || fileParts.isEmpty()) {
                return Response.status(Response.Status.BAD_REQUEST)
                        .entity("{\"error\": \"No file uploaded\"}")
                        .build();
            }

            InputPart filePart = fileParts.get(0);
            InputStream uploadedInputStream = filePart.getBody(InputStream.class, null);

            List<InputPart> filenameParts = uploadForm.get("filename");
            String filename = null;
            if (filenameParts != null && !filenameParts.isEmpty()) {
                filename = filenameParts.get(0).getBody(String.class, null);
            }

            if (filename == null) {
                // Try to get filename from content disposition header
                String contentDisposition = filePart.getHeaders().getFirst("Content-Disposition");
                if (contentDisposition != null) {
                    filename = contentDisposition.substring(contentDisposition.indexOf("filename=") + 10,
                            contentDisposition.length() - 1);
                }
            }

            // Validate file extension
            if (filename == null || !filename.toLowerCase().endsWith(".jar")) {
                return Response.status(Response.Status.BAD_REQUEST)
                        .entity("{\"error\": \"Only JAR files are allowed\"}")
                        .build();
            }

            // Create plugins directory if it doesn't exist
            java.nio.file.Path pluginsDir = Paths.get(pluginManager.getPluginsDirectory());
            if (!Files.exists(pluginsDir)) {
                Files.createDirectories(pluginsDir);
            }

            // Save the uploaded file
            java.nio.file.Path targetPath = pluginsDir.resolve(filename);
            if (Files.exists(targetPath)) {
                return Response.status(Response.Status.CONFLICT)
                        .entity("{\"error\": \"Plugin file already exists\"}")
                        .build();
            }

            try (FileOutputStream outputStream = new FileOutputStream(targetPath.toFile())) {
                byte[] buffer = new byte[8192];
                int bytesRead;
                while ((bytesRead = uploadedInputStream.read(buffer)) != -1) {
                    outputStream.write(buffer, 0, bytesRead);
                }
            }

            // Attempt to load the plugin
            boolean loaded = pluginManager.loadPlugin(targetPath);
            if (loaded) {
                Log.infof("Plugin uploaded and loaded successfully: %s", filename);
                return Response.ok()
                        .entity("{\"message\": \"Plugin uploaded and loaded successfully\", \"filename\": \"" + filename + "\"}")
                        .build();
            } else {
                Log.warnf("Plugin uploaded but failed to load: %s", filename);
                return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                        .entity("{\"error\": \"Plugin uploaded but failed to load\", \"filename\": \"" + filename + "\"}")
                        .build();
            }

        } catch (Exception e) {
            Log.errorf("Error uploading plugin: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to upload plugin: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @GET
    @jakarta.ws.rs.Path("/")
    public Response getAllPlugins() {
        try {
            List<PluginInfo> plugins = pluginManager.getAllPlugins();
            StringBuilder response = new StringBuilder("{\"plugins\": [");
            
            for (int i = 0; i < plugins.size(); i++) {
                PluginInfo plugin = plugins.get(i);
                response.append("{")
                        .append("\"name\":\"").append(plugin.getName()).append("\",")
                        .append("\"version\":\"").append(plugin.getVersion()).append("\",")
                        .append("\"fileName\":\"").append(plugin.getFileName()).append("\",")
                        .append("\"enabled\":").append(plugin.isEnabled())
                        .append("}");
                
                if (i < plugins.size() - 1) {
                    response.append(",");
                }
            }
            
            response.append("]}");
            
            return Response.ok(response.toString()).build();
        } catch (Exception e) {
            Log.errorf("Error retrieving plugins: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to retrieve plugins: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @GET
    @jakarta.ws.rs.Path("/{fileName}")
    public Response getPlugin(@jakarta.ws.rs.PathParam("fileName") String fileName) {
        try {
            PluginInfo plugin = pluginManager.getPlugin(fileName);
            if (plugin == null) {
                return Response.status(Response.Status.NOT_FOUND)
                        .entity("{\"error\": \"Plugin not found: " + fileName + "\"}")
                        .build();
            }

            String response = "{"
                    + "\"name\":\"" + plugin.getName() + "\","
                    + "\"version\":\"" + plugin.getVersion() + "\","
                    + "\"fileName\":\"" + plugin.getFileName() + "\","
                    + "\"filePath\":\"" + plugin.getFilePath() + "\","
                    + "\"enabled\":" + plugin.isEnabled()
                    + "}";

            return Response.ok(response).build();
        } catch (Exception e) {
            Log.errorf("Error retrieving plugin: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to retrieve plugin: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @DELETE
    @jakarta.ws.rs.Path("/{fileName}")
    public Response deletePlugin(@jakarta.ws.rs.PathParam("fileName") String fileName) {
        try {
            // First unload the plugin if it's loaded
            pluginManager.unloadPlugin(fileName);

            // Then delete the file
            java.nio.file.Path pluginsDir = Paths.get(pluginManager.getPluginsDirectory());
            java.nio.file.Path pluginPath = pluginsDir.resolve(fileName);

            if (!Files.exists(pluginPath)) {
                return Response.status(Response.Status.NOT_FOUND)
                        .entity("{\"error\": \"Plugin file not found: " + fileName + "\"}")
                        .build();
            }

            Files.delete(pluginPath);
            Log.infof("Plugin deleted: %s", fileName);

            return Response.ok()
                    .entity("{\"message\": \"Plugin deleted successfully\", \"filename\": \"" + fileName + "\"}")
                    .build();
        } catch (Exception e) {
            Log.errorf("Error deleting plugin: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to delete plugin: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @PUT
    @jakarta.ws.rs.Path("/{fileName}/enable")
    public Response enablePlugin(@jakarta.ws.rs.PathParam("fileName") String fileName) {
        try {
            boolean result = pluginManager.enablePlugin(fileName);
            if (result) {
                return Response.ok()
                        .entity("{\"message\": \"Plugin enabled successfully\", \"filename\": \"" + fileName + "\"}")
                        .build();
            } else {
                return Response.status(Response.Status.NOT_FOUND)
                        .entity("{\"error\": \"Plugin not found: " + fileName + "\"}")
                        .build();
            }
        } catch (Exception e) {
            Log.errorf("Error enabling plugin: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to enable plugin: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @PUT
    @jakarta.ws.rs.Path("/{fileName}/disable")
    public Response disablePlugin(@jakarta.ws.rs.PathParam("fileName") String fileName) {
        try {
            boolean result = pluginManager.disablePlugin(fileName);
            if (result) {
                return Response.ok()
                        .entity("{\"message\": \"Plugin disabled successfully\", \"filename\": \"" + fileName + "\"}")
                        .build();
            } else {
                return Response.status(Response.Status.NOT_FOUND)
                        .entity("{\"error\": \"Plugin not found: " + fileName + "\"}")
                        .build();
            }
        } catch (Exception e) {
            Log.errorf("Error disabling plugin: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to disable plugin: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @POST
    @Path("/refresh")
    public Response refreshPlugins() {
        try {
            pluginManager.refreshPlugins();
            return Response.ok()
                    .entity("{\"message\": \"Plugin refresh completed\"}")
                    .build();
        } catch (Exception e) {
            Log.errorf("Error refreshing plugins: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to refresh plugins: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @GET
    @jakarta.ws.rs.Path("/{fileName}/config")
    public Response getPluginConfig(@jakarta.ws.rs.PathParam("fileName") String fileName) {
        try {
            // Extract plugin name from filename (remove .jar extension)
            String pluginName = fileName.replace(".jar", "");
            Properties config = pluginConfigService.loadPluginConfig(pluginName);

            StringBuilder response = new StringBuilder("{\"config\": {");
            boolean first = true;
            for (String key : config.stringPropertyNames()) {
                if (!first) {
                    response.append(",");
                }
                response.append("\"").append(key).append("\":\"").append(config.getProperty(key)).append("\"");
                first = false;
            }
            response.append("}}");

            return Response.ok(response.toString()).build();
        } catch (Exception e) {
            Log.errorf("Error retrieving plugin config: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to retrieve plugin config: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @POST
    @jakarta.ws.rs.Path("/{fileName}/config")
    @Consumes(MediaType.APPLICATION_JSON)
    public Response updatePluginConfig(@jakarta.ws.rs.PathParam("fileName") String fileName, String configJson) {
        try {
            // Extract plugin name from filename (remove .jar extension)
            String pluginName = fileName.replace(".jar", "");

            // Parse the JSON config (simplified - in real implementation you'd use a proper JSON parser)
            // For now, we'll simulate updating properties
            // In a real implementation, you'd parse the JSON and update individual properties
            // This is a simplified version that creates a new config based on the JSON

            // For demonstration purposes, let's assume the JSON is in the format {"key1": "value1", "key2": "value2"}
            // In a real implementation, you'd use Jackson or similar to parse the JSON
            Properties config = pluginConfigService.loadPluginConfig(pluginName);

            // This is a simplified approach - in reality you'd parse the JSON properly
            // For now, let's just create a default config to simulate
            config.setProperty("updated-at", String.valueOf(System.currentTimeMillis()));

            boolean success = pluginConfigService.savePluginConfig(pluginName, config);
            if (success) {
                return Response.ok()
                        .entity("{\"message\": \"Plugin configuration updated successfully\", \"filename\": \"" + fileName + "\"}")
                        .build();
            } else {
                return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                        .entity("{\"error\": \"Failed to update plugin configuration\"}")
                        .build();
            }
        } catch (Exception e) {
            Log.errorf("Error updating plugin config: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to update plugin config: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @PUT
    @jakarta.ws.rs.Path("/{fileName}/config/{key}")
    @Consumes(MediaType.TEXT_PLAIN)
    public Response updatePluginConfigProperty(@jakarta.ws.rs.PathParam("fileName") String fileName,
                                              @jakarta.ws.rs.PathParam("key") String key,
                                              String value) {
        try {
            // Extract plugin name from filename (remove .jar extension)
            String pluginName = fileName.replace(".jar", "");

            boolean success = pluginConfigService.updatePluginConfigProperty(pluginName, key, value);
            if (success) {
                return Response.ok()
                        .entity("{\"message\": \"Plugin configuration property updated\", \"filename\": \"" + fileName + "\", \"key\": \"" + key + "\", \"value\": \"" + value + "\"}")
                        .build();
            } else {
                return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                        .entity("{\"error\": \"Failed to update plugin configuration property\"}")
                        .build();
            }
        } catch (Exception e) {
            Log.errorf("Error updating plugin config property: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to update plugin config property: " + e.getMessage() + "\"}")
                    .build();
        }
    }

    @DELETE
    @jakarta.ws.rs.Path("/{fileName}/config/{key}")
    public Response removePluginConfigProperty(@jakarta.ws.rs.PathParam("fileName") String fileName,
                                              @jakarta.ws.rs.PathParam("key") String key) {
        try {
            // Extract plugin name from filename (remove .jar extension)
            String pluginName = fileName.replace(".jar", "");

            boolean success = pluginConfigService.removePluginConfigProperty(pluginName, key);
            if (success) {
                return Response.ok()
                        .entity("{\"message\": \"Plugin configuration property removed\", \"filename\": \"" + fileName + "\", \"key\": \"" + key + "\"}")
                        .build();
            } else {
                return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                        .entity("{\"error\": \"Failed to remove plugin configuration property\"}")
                        .build();
            }
        } catch (Exception e) {
            Log.errorf("Error removing plugin config property: %s", e.getMessage());
            return Response.status(Response.Status.INTERNAL_SERVER_ERROR)
                    .entity("{\"error\": \"Failed to remove plugin config property: " + e.getMessage() + "\"}")
                    .build();
        }
    }
}
================================================================================

================================================================================
tech/kayys/silat/runtime/workflow/RuntimeWorkflowDefinitionService.java
Size: 1.8 KB | Modified: 2026-01-15 12:52:37
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


=== SUMMARY ===
Files processed: 12
Directories scanned: 10
Total input size: 48.2 KB
Output size: 52.4 KB
Processing time: 0.00 seconds
