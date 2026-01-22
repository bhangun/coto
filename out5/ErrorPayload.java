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

public class ErrorPayload {


    @NotBlank
    private final String type;

    @NotBlank
    private final String message;

    private final Map<String, Object> details;
    private final boolean retryable;
    private final String originNode;
    private final String originRunId;
    private final int attempt;
    private final int maxAttempts;

    @NotNull
    private final Instant timestamp;

    private final String suggestedAction;
    private final String provenanceRef;

    @JsonCreator
    public ErrorPayload(
        @JsonProperty("type") String type,
        @JsonProperty("message") String message,
        @JsonProperty("details") Map<String, Object> details,
        @JsonProperty("retryable") boolean retryable,
        @JsonProperty("originNode") String originNode,
        @JsonProperty("originRunId") String originRunId,
        @JsonProperty("attempt") int attempt,
        @JsonProperty("maxAttempts") int maxAttempts,
        @JsonProperty("timestamp") Instant timestamp,
        @JsonProperty("suggestedAction") String suggestedAction,
        @JsonProperty("provenanceRef") String provenanceRef
    ) {
        this.type = Objects.requireNonNull(type, "type");
        this.message = Objects.requireNonNull(message, "message");
        this.details = details != null
            ? Collections.unmodifiableMap(new HashMap<>(details))
            : Collections.emptyMap();
        this.retryable = retryable;
        this.originNode = originNode;
        this.originRunId = originRunId;
        this.attempt = attempt;
        this.maxAttempts = maxAttempts;
        this.timestamp = timestamp != null ? timestamp : Instant.now();
        this.suggestedAction = suggestedAction;
        this.provenanceRef = provenanceRef;
    
}