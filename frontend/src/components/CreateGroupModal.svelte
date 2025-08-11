<script>
    import { createEventDispatcher } from "svelte";

    export let isCreating = false;
    export let error = "";
    export let newGroupName = "";
    export let newGroupDescription = "";

    const dispatch = createEventDispatcher();

    function handleInput(event) {
        const { id, value } = event.target;
        dispatch("updateField", { id, value });
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
            <h2>Create New Group</h2>
            <button class="close-button" on:click={handleClose}>
                &times;
            </button>
        </div>

        {#if error}
            <div class="error-message">
                {error}
            </div>
        {/if}

        <form on:submit={handleSubmit}>
            <div class="form-row">
                <label for="group-name">Group Name:</label>
                <input
                    type="text"
                    id="group-name"
                    required
                    bind:value={newGroupName}
                    on:input={handleInput}
                    disabled={isCreating}
                    placeholder="Enter group name"
                />
            </div>

            <div class="form-row">
                <label for="group-description">Description (optional):</label>
                <textarea
                    id="group-description"
                    bind:value={newGroupDescription}
                    on:input={handleInput}
                    disabled={isCreating}
                    placeholder="Describe your group..."
                    rows="3"
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
                    {isCreating ? "Creating..." : "Create Group"}
                </button>
            </div>
        </form>
    </div>
</div>
