<script>
    import { createEventDispatcher } from "svelte";

    export let newPostContent = "";
    export let error = "";
    export let isCreating = false;

    const dispatch = createEventDispatcher();

    function handleInput(event) {
        dispatch("updateContent", event.target.value);
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
            <h2>New Prayer Request</h2>
            <button class="close-button" on:click={handleClose}> × </button>
        </div>

        {#if error}
            <div class="error-message">{error}</div>
        {/if}

        <form on:submit={handleSubmit}>
            <div class="form-row">
                <label for="post-content" class="visually-hidden"
                    >Prayer Request:</label
                >
                <textarea
                    id="post-content"
                    required
                    value={newPostContent}
                    on:input={handleInput}
                    disabled={isCreating}
                    placeholder="Share your prayer request with the group..."
                    rows="5"
                ></textarea>
            </div>

            <div class="modal-actions">
                <button
                    type="button"
                    on:click={handleClose}
                    disabled={isCreating}
                >
                    Cancel
                </button>
                <button type="submit" disabled={isCreating}>
                    {isCreating ? "Posting..." : "Post Request"}
                </button>
            </div>
        </form>
    </div>
</div>
