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

public class DefaultExecutionStateMachine {


    private static final Logger LOG = Logger.getLogger(DefaultExecutionStateMachine.class);

    // Valid transitions map for validation
    private static final Map<ExecutionStatus, Set<ExecutionStatus>> ALLOWED_TRANSITIONS = Map.of(
        ExecutionStatus.CREATED, Set.of(
            ExecutionStatus.RUNNING, 
            ExecutionStatus.CANCELLED
        ),
        ExecutionStatus.RUNNING, Set.of(
            ExecutionStatus.WAITING, 
            ExecutionStatus.RETRYING,
            ExecutionStatus.COMPLETED, 
            ExecutionStatus.FAILED,
            ExecutionStatus.SUSPENDED, 
            ExecutionStatus.CANCELLED
        ),
        ExecutionStatus.WAITING, Set.of(
            ExecutionStatus.RUNNING, 
            ExecutionStatus.FAILED,
            ExecutionStatus.CANCELLED
        ),
        ExecutionStatus.SUSPENDED, Set.of(
            ExecutionStatus.RUNNING, 
            ExecutionStatus.CANCELLED
        ),
        ExecutionStatus.RETRYING, Set.of(
            ExecutionStatus.RUNNING, 
            ExecutionStatus.FAILED
        )
    );

    @Override
    public ExecutionStatus next(ExecutionStatus current, ExecutionSignal signal) {
        ExecutionStatus nextState = computeNextState(current, signal);
        
        if (!isTransitionAllowed(current, nextState)) {
            throw new IllegalStateTransitionException(
                String.format(
                    "Invalid transition from %s to %s via signal %s",
                    current, nextState, signal
                )
            );
        
}