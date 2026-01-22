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

public class InferenceResponse {


    @NotBlank
    private final String requestId;

    @NotBlank
    private final String content;

    private final String model;
    private final int tokensUsed;
    private final long durationMs;
    
    @NotNull
    private final Instant timestamp;

    private final Map<String, Object> metadata;

    @JsonCreator
    public InferenceResponse(
        @JsonProperty("requestId") String requestId,
        @JsonProperty("content") String content,
        @JsonProperty("model") String model,
        @JsonProperty("tokensUsed") int tokensUsed,
        @JsonProperty("durationMs") long durationMs,
        @JsonProperty("timestamp") Instant timestamp,
        @JsonProperty("metadata") Map<String, Object> metadata
    ) {
        this.requestId = Objects.requireNonNull(requestId, "requestId");
        this.content = Objects.requireNonNull(content, "content");
        this.model = model;
        this.tokensUsed = tokensUsed;
        this.durationMs = durationMs;
        this.timestamp = timestamp != null ? timestamp : Instant.now();
        this.metadata = metadata != null
            ? Collections.unmodifiableMap(new HashMap<>(metadata))
            : Collections.emptyMap();
    
}