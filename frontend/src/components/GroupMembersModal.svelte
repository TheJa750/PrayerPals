<script>
    import { createEventDispatcher } from "svelte";

    export let members = [];
    export let error = "";
    export let isLoading = false;
    export let inviteCode = "";
    export let isAdmin = false;

    const dispatch = createEventDispatcher();

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
    <!-- svelte-ignore a11y-click-events-have-key-events -->
    <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
    <div
        class="modal-content text-center"
        on:click|stopPropagation
        role="document"
    >
        <div class="modal-header">
            <h2>{inviteCode}</h2>
            <h2>Group Members ({members.length})</h2>
            <button class="close-button ml-close" on:click={handleClose}>
                &times;
            </button>
        </div>

        {#if error}
            <div class="error-message">{error}</div>
        {/if}

        {#if isLoading}
            <p>Loading members...</p>
        {:else if members.length === 0}
            <p>No members in this group.</p>
        {:else}
            <ul class="member-list">
                {#each members as member}
                    {#if member.role !== "member"}
                        <li>{member.username} - {member.role}</li>
                    {:else if isAdmin}
                        <li>
                            {member.username}
                            <button
                                class="close-button"
                                on:click={() => dispatch("remove", member)}
                            >
                                &times;
                            </button>
                        </li>
                    {:else}
                        <li>{member.username}</li>
                    {/if}
                {/each}
            </ul>
        {/if}
    </div>
</div>
