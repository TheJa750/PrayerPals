<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import { REFRESH_ERROR_MESSAGE } from "../lib/api";
    import { validatePassword, validateUsername } from "../lib/validation";

    export let navigate;

    // User data
    let userGroups = [];
    let isLoadingGroups = true;
    let errorMessage = "";

    // Modal states
    let showCreateGroupModal = false;
    let showJoinGroupModal = false;
    let showChangePasswordModal = false;
    let showChangeUsernameModal = false;

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

    // Change user data
    let newUsername = "";
    let newPassword = "";
    let isChangingUsername = false;
    let isChangingPassword = false;
    let changeError = "";

    // Validation states
    let newUsernameValidation = { isValid: true, errors: [] };
    let newPasswordValidation = { isValid: true, errors: [] };

    // Reactive validation checks
    $: newUsernameValidation = newUsername
        ? validateUsername(newUsername)
        : { isValid: true, errors: [] };
    $: newPasswordValidation = newPassword
        ? validatePassword(newPassword)
        : { isValid: true, errors: [] };

    // Load user groups from the API
    async function loadUserGroups() {
        try {
            isLoadingGroups = true;
            errorMessage = "";

            const groups = await apiRequest("/groups", "GET");
            userGroups = groups || [];
        } catch (error) {
            console.error("Error loading user groups:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

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

    function openChangeUsernameModal() {
        showChangeUsernameModal = true;
    }

    function closeChangeUsernameModal() {
        showChangeUsernameModal = false;
        newUsername = "";
        changeError = "";
    }

    function openChangePasswordModal() {
        showChangePasswordModal = true;
    }

    function closeChangePasswordModal() {
        showChangePasswordModal = false;
        newPassword = "";
        changeError = "";
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

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

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
            const response = await apiRequest(
                `/groups/invite/${inviteCode}`,
                "GET",
            );
            groupPreview = response;
        } catch (error) {
            console.error("Error previewing group:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

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
                `/groups/invite/${inviteCode}/join`,
                "POST",
            );

            closeJoinGroupModal();
            await loadUserGroups();
        } catch (error) {
            console.error("Error joining group:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            joinGroupError =
                error.message || "Failed to join group. Please try again.";
        } finally {
            isJoiningGroup = false;
        }
    }

    async function handleChangeUsername(event) {
        event.preventDefault();
        isChangingUsername = true;
        changeError = "";

        try {
            const response = await apiRequest("/users/update", "PUT", {
                username: newUsername,
            });

            closeChangeUsernameModal();
            await loadUserGroups();
        } catch (error) {
            console.error("Error changing username:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            changeError =
                error.message || "Failed to change username. Please try again.";
        } finally {
            isChangingUsername = false;
        }
    }

    async function handleChangePassword(event) {
        event.preventDefault();
        isChangingPassword = true;
        changeError = "";

        try {
            const response = await apiRequest("/users/update", "PUT", {
                password: newPassword,
            });

            closeChangePasswordModal();

            const response1 = await apiRequest("/logout", "POST");

            navigate("login");
        } catch (error) {
            console.error("Error changing password:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            changeError =
                error.message || "Failed to change password. Please try again.";
        } finally {
            isChangingPassword = false;
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
                <button class="action-button" on:click={openChangeUsernameModal}
                    >Change Username</button
                >
                <button class="action-button" on:click={openChangePasswordModal}
                    >Change Password</button
                >
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
                            <div class="error-message text-center">
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

    <!-- Change Username/Password Modal -->
    {#if showChangeUsernameModal}
        <div
            class="modal-overlay"
            on:click={() => {
                showChangeUsernameModal = false;
            }}
            on:keydown={(e) =>
                e.key === "Escape" && (showChangeUsernameModal = false)}
            role="dialog"
            aria-modal="true"
            tabindex="-1"
        >
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
            <div class="modal-content" on:click|stopPropagation role="document">
                <div class="modal-header">
                    <h2>Change Username</h2>
                    <button
                        class="close-button"
                        on:click={() => {
                            showChangeUsernameModal = false;
                        }}
                    >
                        &times;
                    </button>
                </div>
                <form on:submit={handleChangeUsername}>
                    <div class="form-row">
                        <label for="new-username">New Username:</label>
                        <input
                            type="text"
                            id="new-username"
                            bind:value={newUsername}
                            required
                            placeholder="Enter new username"
                        />
                    </div>
                    {#if !newUsernameValidation.isValid && newUsername}
                        <div class="validation-container">
                            {#each newUsernameValidation.errors as error}
                                <p class="validation-error">{error}</p>
                            {/each}
                        </div>
                    {/if}

                    {#if changeError}
                        <div class="error-message">
                            {changeError}
                        </div>
                    {/if}

                    <div class="modal-actions">
                        <button
                            type="button"
                            on:click={() => {
                                showChangeUsernameModal = false;
                            }}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={isChangingUsername ||
                                !newUsernameValidation.isValid ||
                                !newUsername}
                        >
                            {isChangingUsername
                                ? "Changing..."
                                : "Change Username"}
                        </button>
                    </div>
                </form>
            </div>
        </div>

        <!-- Change Password Modal -->
    {:else if showChangePasswordModal}
        <div
            class="modal-overlay"
            on:click={() => {
                showChangePasswordModal = false;
            }}
            on:keydown={(e) =>
                e.key === "Escape" && (showChangePasswordModal = false)}
            role="dialog"
            aria-modal="true"
            tabindex="-1"
        >
            <!-- svelte-ignore a11y_click_events_have_key_events -->
            <!-- svelte-ignore a11y_no_noninteractive_element_interactions -->
            <div class="modal-content" on:click|stopPropagation role="document">
                <div class="modal-header">
                    <h2>Change Password</h2>
                    <button
                        class="close-button"
                        on:click={() => {
                            showChangePasswordModal = false;
                        }}
                    >
                        &times;
                    </button>
                </div>
                <form on:submit={handleChangePassword}>
                    <div class="form-row">
                        <label for="new-password">New Password:</label>
                        <input
                            type="password"
                            id="new-password"
                            bind:value={newPassword}
                            required
                            placeholder="Enter new password"
                        />
                    </div>
                    {#if !newPasswordValidation.isValid && newPassword}
                        <div class="validation-container">
                            {#each newPasswordValidation.errors as error}
                                <p class="validation-error">{error}</p>
                            {/each}
                        </div>
                    {/if}

                    {#if changeError}
                        <div class="error-message">
                            {changeError}
                        </div>
                    {/if}

                    <div class="modal-actions">
                        <button
                            type="button"
                            on:click={() => {
                                showChangePasswordModal = false;
                            }}
                        >
                            Cancel
                        </button>
                        <button
                            type="submit"
                            disabled={isChangingPassword ||
                                !newPasswordValidation.isValid ||
                                !newPassword}
                        >
                            {isChangingPassword
                                ? "Changing..."
                                : "Change Password"}
                        </button>
                    </div>
                </form>
            </div>
        </div>
    {/if}
</div>
