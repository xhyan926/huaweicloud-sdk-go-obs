/**
 * Copyright 2019 Huawei Technologies Co.,Ltd.
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License.  You may obtain a copy of the
 * License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed
 * under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
 * CONDITIONS OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package com.obs.services.model;

import static org.junit.Assert.assertEquals;
import static org.junit.Assert.assertFalse;
import static org.junit.Assert.assertNotNull;
import static org.junit.Assert.assertTrue;

import com.obs.aitool.AIGenerated;
import com.obs.services.internal.utils.CallCancelHandler;

import java.util.Arrays;
import java.util.Collection;

import org.junit.Test;
import org.junit.runner.RunWith;
import org.junit.runners.Parameterized;
import org.junit.runners.Parameterized.Parameter;
import org.junit.runners.Parameterized.Parameters;

@RunWith(Parameterized.class)
public class ResumableTransferHandleTest {

    @Parameter(0)
    public String testName;

    @Parameter(1)
    public String testCategory;

    @Parameters(name = "{0}")
    public static Collection<Object[]> testData() {
        return Arrays.asList(new Object[][] {
            {"InitialNotPausedNotCancelled_STATE", "INITIAL_STATE"},
            {"PauseSetsPausedState_PAUSE", "PAUSE"},
            {"CancelSetsCancelledState_CANCEL", "CANCEL"},
            {"CancelFromPausedGoesToCancelled_COMBINED", "COMBINED"},
            {"ResumeFromPausedBackToActive_RESUME", "RESUME"},
            {"BindWithNullCreatesHandler_BIND", "BIND"},
            {"BindWithExistingHandler_BIND", "BIND_EXISTING"},
            {"GetCancelHandlerReturnsBound_GET", "GET"},
            {"PauseWhenPausedThrows_ILLEGAL", "PAUSE_WHEN_PAUSED"},
            {"PauseWhenCancelledThrows_ILLEGAL", "PAUSE_WHEN_CANCELLED"},
            {"ResumeWhenActiveThrows_ILLEGAL", "RESUME_WHEN_ACTIVE"},
            {"ResumeWhenCancelledThrows_ILLEGAL", "RESUME_WHEN_CANCELLED"},
            {"CancelIsIdempotent_IDEMPOTENT", "CANCEL_IDEMPOTENT"},
        });
    }

    @AIGenerated(author = "yanliwei", date = "2026-04-16",
            description = "Parameterized test for ResumableTransferHandle state machine transitions")
    @Test
    public void should_verify_correct_state_transitions_when_operations_applied() {
        switch (testCategory) {
            case "INITIAL_STATE":
                runInitialStateTest();
                break;
            case "PAUSE":
                runPauseTest();
                break;
            case "CANCEL":
                runCancelTest();
                break;
            case "COMBINED":
                runCombinedTest();
                break;
            case "RESUME":
                runResumeTest();
                break;
            case "BIND":
                runBindNullTest();
                break;
            case "BIND_EXISTING":
                runBindExistingTest();
                break;
            case "GET":
                runGetCancelHandlerTest();
                break;
            case "PAUSE_WHEN_PAUSED":
                runPauseWhenPausedThrowsTest();
                break;
            case "PAUSE_WHEN_CANCELLED":
                runPauseWhenCancelledThrowsTest();
                break;
            case "RESUME_WHEN_ACTIVE":
                runResumeWhenActiveThrowsTest();
                break;
            case "RESUME_WHEN_CANCELLED":
                runResumeWhenCancelledThrowsTest();
                break;
            case "CANCEL_IDEMPOTENT":
                runCancelIdempotentTest();
                break;
            default:
                break;
        }
    }

    private void runInitialStateTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);

        assertFalse("should not be paused initially", handle.isPaused());
        assertFalse("should not be cancelled initially", handle.isCancelled());
        assertFalse("should not be paused or cancelled initially", handle.isPausedOrCancelled());
    }

    private void runPauseTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);

        handle.pause();
        assertTrue("should be paused after pause()", handle.isPaused());
        assertFalse("should not be cancelled after pause()", handle.isCancelled());
        assertTrue("paused or cancelled should be true after pause()", handle.isPausedOrCancelled());
    }

    private void runCancelTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);

        handle.cancel();
        assertTrue("should be cancelled after cancel()", handle.isCancelled());
        assertFalse("should not be paused after cancel()", handle.isPaused());
        assertTrue("paused or cancelled should be true after cancel()", handle.isPausedOrCancelled());
    }

    private void runCombinedTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);

        handle.pause();
        assertTrue("should be paused after pause()", handle.isPaused());
        assertFalse("should not be cancelled after pause()", handle.isCancelled());

        handle.cancel();
        // cancel 从 PAUSED 进入 CANCELLED 终态，不再处于 PAUSED
        assertFalse("should not be paused after cancel overrode pause", handle.isPaused());
        assertTrue("should be cancelled after cancel()", handle.isCancelled());
        assertTrue("paused or cancelled should be true", handle.isPausedOrCancelled());
    }

    private void runResumeTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);

        handle.pause();
        assertTrue("should be paused after pause()", handle.isPaused());

        handle.resume();
        assertFalse("should not be paused after resume()", handle.isPaused());
        assertFalse("should not be cancelled after resume()", handle.isCancelled());
        assertFalse("paused or cancelled should be false after resume()", handle.isPausedOrCancelled());
    }

    private void runBindNullTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);
        assertNotNull("cancelHandler should not be null after bind(null)", handle.getCancelHandler());
    }

    private void runBindExistingTest() {
        CallCancelHandler existingHandler = new CallCancelHandler();
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(existingHandler);
        assertEquals("should use provided cancelHandler", existingHandler, handle.getCancelHandler());
    }

    private void runGetCancelHandlerTest() {
        CallCancelHandler handler = new CallCancelHandler();
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(handler);
        assertEquals("getCancelHandler should return bound handler", handler, handle.getCancelHandler());
    }

    private void runPauseWhenPausedThrowsTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);
        handle.pause();

        try {
            handle.pause();
            assertTrue("pause() on PAUSED state should throw IllegalStateException", false);
        } catch (IllegalStateException e) {
            assertTrue("exception message should mention current state",
                    e.getMessage().contains("PAUSED"));
        }
    }

    private void runPauseWhenCancelledThrowsTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);
        handle.cancel();

        try {
            handle.pause();
            assertTrue("pause() on CANCELLED state should throw IllegalStateException", false);
        } catch (IllegalStateException e) {
            assertTrue("exception message should mention current state",
                    e.getMessage().contains("CANCELLED"));
        }
    }

    private void runResumeWhenActiveThrowsTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);

        try {
            handle.resume();
            assertTrue("resume() on ACTIVE state should throw IllegalStateException", false);
        } catch (IllegalStateException e) {
            assertTrue("exception message should mention current state",
                    e.getMessage().contains("ACTIVE"));
        }
    }

    private void runResumeWhenCancelledThrowsTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);
        handle.cancel();

        try {
            handle.resume();
            assertTrue("resume() on CANCELLED state should throw IllegalStateException", false);
        } catch (IllegalStateException e) {
            assertTrue("exception message should mention current state",
                    e.getMessage().contains("CANCELLED"));
        }
    }

    private void runCancelIdempotentTest() {
        ResumableTransferHandle handle = new ResumableTransferHandle();
        handle.bind(null);

        handle.cancel();
        assertTrue("should be cancelled after first cancel()", handle.isCancelled());

        // 二次 cancel 应为幂等操作，不抛异常
        handle.cancel();
        assertTrue("should still be cancelled after second cancel()", handle.isCancelled());
    }
}
