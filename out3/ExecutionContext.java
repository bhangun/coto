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

public interface ExecutionContext {


    /**
     * Get engine context (global state)
     */
    EngineContext engine();

    /**
     * Get current execution token (immutable snapshot)
     */
    ExecutionToken token();

    /**
     * Get tenant context
     */
    TenantContext tenantContext();

    /**
     * Update execution status (creates new token)
     */
    void updateStatus(ExecutionStatus status);

    /**
     * Update current phase (creates new token)
     */
    void updatePhase(InferencePhase phase);

    /**
     * Increment retry attempt
     */
    void incrementAttempt();

    /**
     * Get execution variables (mutable view)
     */
    Map<String, Object> variables();

    /**
     * Put variable
     */
    void putVariable(String key, Object value);

    /**
     * Get variable
     */
    <T> Optional<T> getVariable(String key, Class<T> type);

    /**
     * Remove variable
     */
    void removeVariable(String key);

    /**
     * Get metadata
     */
    Map<String, Object> metadata();

    /**
     * Put metadata
     */
    void putMetadata(String key, Object value);

    /**
     * Replace entire token (for state restoration)
     */
    void replaceToken(ExecutionToken newToken);

    /**
     * Check if context has error
     */
    boolean hasError();

    /**
     * Get error if present
     */
    Optional<Throwable> getError();

    /**
     * Set error
     */
    void setError(Throwable error);

    /**
     * Clear error
     */
    void clearError();

}