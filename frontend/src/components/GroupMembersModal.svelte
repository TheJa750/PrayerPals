<script>
    import { createEventDispatcher } from "svelte";

    export let members = [];
    export let error = "";
    export let isLoading = false;
    export let inviteCode = "";
    export let isAdmin = false;

    let copied = false;

    const dispatch = createEventDispatcher();

    function handleClose() {
        dispatch("close");
    }

    function copyInviteCode() {
        navigator.clipboard.writeText(inviteCode);
        copied = true;
        setTimeout(() => (copied = false), 2000);
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
            <div>
                <h3>Invite Code:</h3>
                <h3 class="invite-code" on:click={copyInviteCode}>
                    {inviteCode}
                </h3>
                {#if copied}
                    <div class="copied-popup">Copied!</div>
                {/if}
            </div>
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
            <div class="member-list-container">
                <ul class="member-list">
                    {#each members as member}
                        {#if member.role !== "member"}
                            <li>
                                <div class="member-list-username">
                                    <span class="face-icon">🧑</span>
                                    <span
                                        >{member.username} - ★ {member.role}</span
                                    >
                                </div>
                            </li>
                        {:else if isAdmin}
                            <li class="admin-member-view">
                                <div class="member-list-username">
                                    <span class="face-icon">🧑</span>
                                    <span>{member.username}</span>
                                </div>
                                <div>
                                    <div class="tooltip-container">
                                        <button
                                            class="close-button"
                                            on:click={() =>
                                                dispatch("remove", member)}
                                        >
                                            &times;
                                        </button>
                                        <span class="tooltip-text">
                                            Moderate User
                                        </span>
                                    </div>
                                    <div class="tooltip-container">
                                        <button
                                            class="promote-button"
                                            on:click={() =>
                                                dispatch("promote", member)}
                                        >
                                            ★
                                        </button>
                                        <span class="tooltip-text">
                                            Promote member
                                        </span>
                                    </div>
                                </div>
                            </li>
                        {:else}
                            <li>
                                <div class="member-list-username">
                                    <span class="face-icon">🧑</span>
                                    <span>{member.username}</span>
                                </div>
                            </li>
                        {/if}
                    {/each}
                </ul>
            </div>
        {/if}
    </div>
</div>
