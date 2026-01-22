package tech.kayys.golek.inference.api;

import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import org.jetbrains.annotations.Nullable;
import java.time.Duration;
import java.util.*;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import java.util.Objects;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import java.time.Instant;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import java.util.*;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import java.time.Instant;
import java.util.Collections;
import java.util.HashMap;
import java.util.Map;
import java.util.Objects;
import com.fasterxml.jackson.annotation.JsonCreator;
import com.fasterxml.jackson.annotation.JsonProperty;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.NotNull;
import org.apache.commons.codec.digest.DigestUtils;
import java.time.Instant;
import java.util.*;
import tech.kayys.golek.inference.kernel.pipeline.InferencePhase;
import java.io.Serializable;
import java.time.Instant;
import java.util.Collections;
import java.util.Map;
import java.util.Objects;
import java.util.UUID;
import java.util.concurrent.ConcurrentHashMap;
import jakarta.enterprise.context.ApplicationScoped;
import org.jboss.logging.Logger;
import java.util.Map;
import java.util.Set;
import tech.kayys.golek.inference.api.TenantContext;
import tech.kayys.golek.inference.kernel.engine.EngineContext;
import tech.kayys.golek.inference.kernel.pipeline.InferencePhase;
import java.util.Map;
import java.util.Optional;
import tech.kayys.golek.inference.api.TenantContext;
import tech.kayys.golek.inference.kernel.engine.EngineContext;
import tech.kayys.golek.inference.kernel.pipeline.InferencePhase;
import java.util.Map;
import java.util.Optional;
import java.util.concurrent.atomic.AtomicReference;
import java.util.Arrays;
import java.util.Comparator;
import java.util.List;
import tech.kayys.golek.inference.kernel.execution.ExecutionContext;
import io.smallrye.mutiny.Uni;

public interface ModelRunner {

    
    /**
     * Initialize the runner with model manifest and configuration
     * @param manifest Model metadata and artifact locations
     * @param config Runner-specific configuration
     * @param tenantContext Current tenant context for isolation
     * @throws ModelLoadException if initialization fails
     */
    void initialize(
        ModelManifest manifest, 
        Map<String, Object> config,
        TenantContext tenantContext
    ) throws ModelLoadException;
    
    /**
     * Execute synchronous inference
     * @param request Inference request with inputs
     * @param context Request context with timeout, priority, etc.
     * @return Inference response with outputs
     * @throws InferenceException if execution fails
     */
    InferenceResponse infer(
        InferenceRequest request,
        RequestContext context
    ) throws InferenceException;
    
    /**
     * Execute asynchronous inference with callback
     * @param request Inference request
     * @param context Request context
     * @return CompletionStage for async processing
     */
    CompletionStage<InferenceResponse> inferAsync(
        InferenceRequest request,
        RequestContext context
    );
    
    /**
     * Health check for this runner instance
     * @return Health status with diagnostics
     */
    HealthStatus health();
    
    /**
     * Get current resource utilization metrics
     * @return Resource usage snapshot
     */
    ResourceMetrics getMetrics();
    
    /**
     * Warm up the model (optional optimization)
     * @param sampleInputs Sample inputs for warming
     */
    default void warmup(List<InferenceRequest> sampleInputs) {
        // Default no-op
    
}