<script>
    import { createEventDispatcher } from "svelte";

    export let isDeleting = false;
    export let error = "";

    const dispatch = createEventDispatcher();

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
    <div
        class="modal-content text-center"
        on:click|stopPropagation
        role="document"
    >
        <div class="modal-header">
            <h2>Delete Post</h2>
            <button class="close-button" on:click={handleClose}> × </button>
        </div>

        {#if error}
            <div class="error-message">{error}</div>
        {/if}

        <p>Are you sure you want to delete this post?</p>

        <div class="modal-actions">
            <button type="button" on:click={handleClose} disabled={isDeleting}>
                Cancel
            </button>
            <button type="button" on:click={handleSubmit} disabled={isDeleting}>
                {isDeleting ? "Deleting..." : "Delete"}
            </button>
        </div>
    </div>
</div>
