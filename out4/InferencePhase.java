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

public enum InferencePhase {

    
    /**
     * Phase 1: Pre-validation checks
     * - Request structure validation
     * - Basic sanity checks
     * - Early rejection of malformed requests
     */
    PRE_VALIDATE(1, "Pre-Validation"),
    
    /**
     * Phase 2: Deep validation
     * - Schema validation (JSON Schema)
     * - Content safety checks
     * - Input format verification
     * - Model compatibility checks
     */
    VALIDATE(2, "Validation"),
    
    /**
     * Phase 3: Authorization & access control
     * - Tenant verification
     * - Model access permissions
     * - Feature flag checks
     * - Quota verification
     */
    AUTHORIZE(3, "Authorization"),
    
    /**
     * Phase 4: Intelligent routing & provider selection
     * - Model-to-provider mapping
     * - Multi-factor scoring
     * - Load balancing
     * - Availability checks
     */
    ROUTE(4, "Routing"),
    
    /**
     * Phase 5: Request transformation & enrichment
     * - Prompt templating
     * - Context injection
     * - Parameter normalization
     * - Request mutation
     */
    PRE_PROCESSING(5, "Pre-Processing"),
    
    /**
     * Phase 6: Actual provider dispatch
     * - Provider invocation
     * - Streaming/batch execution
     * - Circuit breaker protection
     * - Fallback handling
     */
    PROVIDER_DISPATCH(6, "Provider Dispatch"),
    
    /**
     * Phase 7: Response post-processing
     * - Output validation
     * - Format normalization
     * - Metadata enrichment
     * - Content moderation
     */
    POST_PROCESSING(7, "Post-Processing"),
    
    /**
     * Phase 8: Audit logging
     * - Event recording
     * - Provenance tracking
     * - Compliance logging
     * - Immutable audit trail
     */
    AUDIT(8, "Audit"),
    
    /**
     * Phase 9: Observability & metrics
     * - Metrics emission
     * - Trace completion
     * - Performance tracking
     * - Cost accounting
     */
    OBSERVABILITY(9, "Observability"),
    
    /**
     * Phase 10: Resource cleanup
     * - Cache invalidation
     * - Connection release
     * - Quota release
     * - Temporary resource cleanup
     */
    CLEANUP(10, "Cleanup");

    private final int order;
    private final String displayName;

    InferencePhase(int order, String displayName) {
        this.order = order;
        this.displayName = displayName;
    
}