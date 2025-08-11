<script>
    // @ts-nocheck

    import { createEventDispatcher } from "svelte";

    export let inviteCode = "";
    export let error = "";
    export let isLoadingPreview = false;
    export let groupPreview = null;
    export let isJoiningGroup = false;

    const dispatch = createEventDispatcher();

    function handleInput(event) {
        dispatch("updateCode", event.target.value);
    }

    function handlePreview(event) {
        event.preventDefault();
        dispatch("preview");
    }

    function handleJoin() {
        dispatch("join");
    }

    function handleClose() {
        dispatch("close");
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
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
    <div class="modal-content" on:click|stopPropagation role="document">
        <!-- Before group preview -->
        {#if !groupPreview}
            <div class="modal-header">
                <h2>Enter Invite Code</h2>
                <button class="close-button" on:click={handleClose}
                    >&times;</button
                >
            </div>

            <form on:submit={handlePreview}>
                <div class="form-row">
                    <label for="invite-code">Invite Code:</label>
                    <input
                        type="text"
                        id="invite-code"
                        value={inviteCode}
                        on:input={handleInput}
                        required
                        placeholder="Enter invite code"
                        disabled={isLoadingPreview}
                    />
                </div>
                {#if error}
                    <div class="error-message">{error}</div>
                {/if}
                <div class="modal-actions">
                    <button
                        type="button"
                        on:click={handleClose}
                        disabled={isLoadingPreview}>Cancel</button
                    >
                    <button type="submit" disabled={isLoadingPreview}>
                        {isLoadingPreview ? "Loading..." : "Find Group"}
                    </button>
                </div>
            </form>
        {:else}
            <div class="modal-header">
                <h2>Join Group</h2>
                <button class="close-button" on:click={handleClose}
                    >&times;</button
                >
            </div>
            <div class="group-card clickable" role="button" tabindex="0">
                <h3>{groupPreview.name}</h3>
                {#if groupPreview.description}
                    <p class="group-description">{groupPreview.description}</p>
                {/if}
            </div>
            {#if error}
                <div class="error-message">{error}</div>
            {/if}
            <div class="modal-actions">
                <button
                    type="button"
                    on:click={handleClose}
                    disabled={isJoiningGroup}>Cancel</button
                >
                <button
                    type="button"
                    on:click={handleJoin}
                    disabled={isJoiningGroup}
                >
                    {isJoiningGroup ? "Joining..." : "Join"}
                </button>
            </div>
        {/if}
    </div>
</div>
