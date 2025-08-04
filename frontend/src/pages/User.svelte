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
    let isLoadingPreview = false;
    let groupPreview = null;

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
        groupPreview = null;
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

    async function handlePreviewGroup(event) {
        event.preventDefault();
        isLoadingPreview = true;
        joinGroupError = "";
        groupPreview = null;

        try {
            const response = await apiRequest(`/groups/${inviteCode}`, "GET");
            groupPreview = response;
        } catch (error) {
            console.error("Error previewing group:", error);
            joinGroupError =
                error.message || "Failed to find group. Please check the code.";
        } finally {
            isLoadingPreview = false;
        }
    }

    async function handleJoinGroup() {
        if (!groupPreview || isJoiningGroup) return;

        isJoiningGroup = true;
        joinGroupError = "";

        try {
            const joinResponse = await apiRequest(
                `/groups/${inviteCode}/join`,
                "POST",
            );

            closeJoinGroupModal();
            await loadUserGroups();
        } catch (error) {
            console.error("Error joining group:", error);
            joinGroupError =
                error.message || "Failed to join group. Please try again.";
        } finally {
            isJoiningGroup = false;
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
                            disabled={isLoadingPreview}
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

    <!-- Join Group Modal -->
    {#if showJoinGroupModal}
        <div
            class="modal-overlay"
            on:click={closeJoinGroupModal}
            on:keydown={(e) => e.key === "Escape" && closeJoinGroupModal()}
            role="dialog"
            aria-modal="true"
            tabindex="-1"
        >
            {#if !groupPreview}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
                <div
                    class="modal-content"
                    on:click|stopPropagation
                    role="document"
                >
                    <div class="modal-header">
                        <h2>Enter Invite Code</h2>
                        <button
                            class="close-button"
                            on:click={closeJoinGroupModal}
                        >
                            &times;
                        </button>
                    </div>

                    <form on:submit={handlePreviewGroup}>
                        <div class="form-row">
                            <label for="invite-code">Invite Code:</label>
                            <input
                                type="text"
                                id="invite-code"
                                bind:value={inviteCode}
                                required
                                placeholder="Enter invite code"
                            />
                        </div>

                        {#if joinGroupError}
                            <div class="error-message">
                                {joinGroupError}
                            </div>
                        {/if}

                        <div class="modal-actions">
                            <button
                                type="button"
                                on:click={closeJoinGroupModal}
                                disabled={isLoadingPreview}
                            >
                                Cancel
                            </button>
                            <button type="submit" disabled={isLoadingPreview}>
                                {isLoadingPreview ? "Loading..." : "Find Group"}
                            </button>
                        </div>
                    </form>
                </div>
            {:else}
                <!-- svelte-ignore a11y-click-events-have-key-events -->
                <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
                <div
                    class="modal-content"
                    on:click|stopPropagation
                    role="document"
                >
                    <div class="modal-header">
                        <h2>Join Group</h2>
                        <button
                            class="close-button"
                            on:click={closeJoinGroupModal}
                        >
                            &times;
                        </button>
                    </div>

                    <div
                        class="group-card clickable"
                        on:click={handleJoinGroup}
                        role="button"
                        tabindex="0"
                    >
                        <h3>{groupPreview.name}</h3>
                        {#if groupPreview.description}
                            <p class="group-description">
                                {groupPreview.description}
                            </p>
                        {/if}
                    </div>

                    {#if joinGroupError}
                        <div class="error-message">
                            {joinGroupError}
                        </div>
                    {/if}

                    <div class="modal-actions">
                        <button
                            type="button"
                            on:click={closeJoinGroupModal}
                            disabled={isJoiningGroup}
                        >
                            Cancel
                        </button>
                        <button
                            type="button"
                            disabled={isJoiningGroup}
                            on:click={handleJoinGroup}
                        >
                            Join
                        </button>
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>
