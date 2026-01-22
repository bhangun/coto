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

public class AuditPayload {


    @NotNull
    private final Instant timestamp;

    @NotBlank
    private final String runId;

    private final String nodeId;

    @NotNull
    private final Actor actor;

    @NotBlank
    private final String event;

    @NotBlank
    private final String level;

    private final List<String> tags;
    private final Map<String, Object> metadata;
    private final Map<String, Object> contextSnapshot;

    @NotBlank
    private final String hash;

    @JsonCreator
    public AuditPayload(
        @JsonProperty("timestamp") Instant timestamp,
        @JsonProperty("runId") String runId,
        @JsonProperty("nodeId") String nodeId,
        @JsonProperty("actor") Actor actor,
        @JsonProperty("event") String event,
        @JsonProperty("level") String level,
        @JsonProperty("tags") List<String> tags,
        @JsonProperty("metadata") Map<String, Object> metadata,
        @JsonProperty("contextSnapshot") Map<String, Object> contextSnapshot,
        @JsonProperty("hash") String hash
    ) {
        this.timestamp = timestamp != null ? timestamp : Instant.now();
        this.runId = Objects.requireNonNull(runId, "runId");
        this.nodeId = nodeId;
        this.actor = Objects.requireNonNull(actor, "actor");
        this.event = Objects.requireNonNull(event, "event");
        this.level = Objects.requireNonNull(level, "level");
        this.tags = tags != null
            ? Collections.unmodifiableList(new ArrayList<>(tags))
            : Collections.emptyList();
        this.metadata = metadata != null
            ? Collections.unmodifiableMap(new HashMap<>(metadata))
            : Collections.emptyMap();
        this.contextSnapshot = contextSnapshot != null
            ? Collections.unmodifiableMap(new HashMap<>(contextSnapshot))
            : Collections.emptyMap();
        this.hash = Objects.requireNonNull(hash, "hash");
    
}