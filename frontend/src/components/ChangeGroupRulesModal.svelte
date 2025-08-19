<script>
    import { createEventDispatcher } from "svelte";

    export let groupRules = "";
    export let error = "";
    export let isSaving = false;

    const dispatch = createEventDispatcher();

    function handleInput(event) {
        const { value } = event.target;
        dispatch("updateRules", value);
    }

    function handleClose() {
        dispatch("close");
    }

    function handleSubmit(event) {
        event.preventDefault();
        dispatch("submit");
    }
</script>

<div
    class="modal-overlay"
    on:click={handleClose}
    on:keydown={(e) => e.key === "Escape" && handleClose()}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
>
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <div class="modal-content" on:click|stopPropagation role="document">
        <div class="modal-header">
            <h2>Change Group Rules</h2>
            <button class="close-button" on:click={handleClose}>
                &times;
            </button>
        </div>

        <form on:submit={handleSubmit}>
            <div class="form-row">
                <label for="group-rules">Group Rules:</label>
                <textarea
                    id="group-rules"
                    required
                    bind:value={groupRules}
                    on:input={handleInput}
                    disabled={isSaving}
                    placeholder="Enter group rules"
                    rows="25"
                    maxlength="1500"
                ></textarea>
            </div>

            {#if error}
                <div class="error-message text-center">
                    {error}
                </div>
            {/if}

            <div class="modal-actions">
                <button
                    type="button"
                    on:click={handleClose}
                    disabled={isSaving}
                >
                    Cancel
                </button>
                <button type="submit" disabled={isSaving}>
                    {isSaving ? "Saving..." : "Save Rules"}
                </button>
            </div>
        </form>
    </div>
</div>
