<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import Group from "./Group.svelte";
    export let navigate;

    // User data
    let userGroups = [];
    let isLoadingGroups = true;
    let errorMessage = "";

    // Modal states
    let showCreateGroupModal = false;
    let showJoinGroupModal = false;

    // Create group data
    let newGroupName = "";
    let newGroupDescription = "";
    let isCreatingGroup = false;
    let createGroupError = "";

    // Join group data
    let inviteCode = "";
    let isJoiningGroup = false;
    let joinGroupError = "";

    // Load user groups from the API
    async function loadUserGroups() {
        try {
            isLoadingGroups = true;
            errorMessage = "";

            const groups = await apiRequest("/groups", "GET");
            userGroups = groups || [];
        } catch (error) {
            console.error("Error loading user groups:", error);
            errorMessage = "Failed to load groups. Please try again later.";
        } finally {
            isLoadingGroups = false;
        }
    }

    // Modal functions
    function openCreateGroupModal() {
        showCreateGroupModal = true;
        newGroupName = "";
        newGroupDescription = "";
        createGroupError = "";
    }

    function closeCreateGroupModal() {
        showCreateGroupModal = false;
    }

    function openJoinGroupModal() {
        showJoinGroupModal = true;
        inviteCode = "";
        joinGroupError = "";
    }

    function closeJoinGroupModal() {
        showJoinGroupModal = false;
    }

    async function handleCreateGroup(event) {
        event.preventDefault();
        isCreatingGroup = true;
        createGroupError = "";

        try {
            const groupData = {
                name: newGroupName,
                description: newGroupDescription || "",
            };

            const newGroup = await apiRequest("/groups", "POST", groupData);

            closeCreateGroupModal();
            await loadUserGroups();
        } catch (error) {
            console.error("Error creating group:", error);
            createGroupError =
                error.message || "Failed to create group. Please try again.";
        } finally {
            isCreatingGroup = false;
        }
    }

    onMount(() => {
        loadUserGroups();
    });
</script>

<div class="dashboard-container">
    <h1 class="text-center">Welcome to Prayer Pals!</h1>

    <div class="dashboard-grid">
        <section class="dashboard-section side-section">
            <h2 class="text-center">Actions</h2>
            <div class="actions-grid">
                <button class="action-button" on:click={openCreateGroupModal}
                    >Create New Group</button
                >
                <button class="action-button" on:click={openJoinGroupModal}
                    >Join Group</button
                >
            </div>
        </section>

        <section class="dashboard-section main-section">
            <h2 class="text-center">My Groups</h2>

            {#if isLoadingGroups}
                <p>Loading your groups...</p>
            {:else if errorMessage}
                <div class="error-message">
                    Error loading groups: {errorMessage}
                </div>
            {:else if userGroups.length === 0}
                <p>You are not a member of any groups yet.</p>
            {:else}
                <div class="groups-list">
                    {#each userGroups as group}
                        <div
                            class="group-card clickable"
                            on:click={() => navigate("group", group.id)}
                            on:keydown={(e) =>
                                e.key === "Enter" &&
                                navigate("group", group.id)}
                            tabindex="0"
                            role="button"
                            aria-label="View {group.name} group"
                        >
                            <h3>{group.name}</h3>
                            {#if group.description}
                                <p class="group-description">
                                    {group.description}
                                </p>
                            {/if}
                        </div>
                    {/each}
                </div>
            {/if}
        </section>

        <section class="dashboard-section side-section">
            <h2 class="text-center">Account</h2>
            <div class="actions-grid">
                <button class="action-button">Change Username</button>
                <button class="action-button">Change Password</button>
            </div>
        </section>
    </div>

    <!-- Create Group Modal -->
    {#if showCreateGroupModal}
        <div
            class="modal-overlay"
            on:click={closeCreateGroupModal}
            on:keydown={(e) => e.key === "Escape" && closeCreateGroupModal()}
            role="dialog"
            aria-modal="true"
            tabindex="-1"
        >
            <!-- svelte-ignore a11y-click-events-have-key-events -->
            <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
            <div class="modal-content" on:click|stopPropagation role="document">
                <div class="modal-header">
                    <h2>Create New Group</h2>
                    <button
                        class="close-button"
                        on:click={closeCreateGroupModal}
                    >
                        &times;
                    </button>
                </div>

                {#if createGroupError}
                    <div class="error-message">
                        {createGroupError}
                    </div>
                {/if}

                <form on:submit={handleCreateGroup}>
                    <div class="form-row">
                        <label for="group-name">Group Name:</label>
                        <input
                            type="text"
                            id="group-name"
                            required
                            bind:value={newGroupName}
                            disabled={isCreatingGroup}
                            placeholder="Enter group name"
                        />
                    </div>

                    <div class="form-row">
                        <label for="group-description"
                            >Description (optional):</label
                        >
                        <textarea
                            id="group-description"
                            bind:value={newGroupDescription}
                            disabled={isCreatingGroup}
                            placeholder="Describe your group..."
                            rows="3"
                        ></textarea>
                    </div>

                    <div class="modal-actions">
                        <button
                            type="button"
                            on:click={closeCreateGroupModal}
                            disabled={isCreatingGroup}
                        >
                            Cancel
                        </button>
                        <button type="submit" disabled={isCreatingGroup}>
                            {isCreatingGroup ? "Creating..." : "Create Group"}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    {/if}
</div>
