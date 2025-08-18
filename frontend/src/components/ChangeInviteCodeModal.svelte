<script>
    import { createEventDispatcher } from "svelte";

    export let inviteCode = "";
    export let error = "";
    export let isChangingCode = false;

    const dispatch = createEventDispatcher();

    function handleInput(event) {
        const { value } = event.target;
        dispatch("updateCode", value);
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
            <h2>Change Invitation Code</h2>
            <button class="close-button" on:click={handleClose}>
                &times;
            </button>
        </div>

        <form on:submit={handleSubmit}>
            <div class="form-row">
                <label for="invite-code">New Invite Code:</label>
                <input
                    type="text"
                    id="invite-code"
                    required
                    bind:value={inviteCode}
                    on:input={handleInput}
                    disabled={isChangingCode}
                    placeholder="Enter new invite code"
                    maxlength="6"
                />
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
                    disabled={isChangingCode}
                >
                    Cancel
                </button>
                <button type="submit" disabled={isChangingCode}>
                    {isChangingCode ? "Changing..." : "Change Code"}
                </button>
            </div>
        </form>
    </div>
</div>
