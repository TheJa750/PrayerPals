<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import { REFRESH_ERROR_MESSAGE } from "../lib/api";
    import { validatePassword, validateUsername } from "../lib/validation";

    import CreateGroupModal from "../components/CreateGroupModal.svelte";
    import JoinGroupModal from "../components/JoinGroupModal.svelte";

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

    function handleModalFieldUpdate(event) {
        const { id, value } = event.detail;
        if (id === "group-name") {
            newGroupName = value;
        } else if (id === "group-description") {
            newGroupDescription = value;
        }
    }

    async function handleCreateGroup(event) {
        event.preventDefault();
        console.log("Creating group:", newGroupName, newGroupDescription);
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
        <CreateGroupModal
            on:close={closeCreateGroupModal}
            on:updateField={handleModalFieldUpdate}
            on:submit={handleCreateGroup}
            {newGroupName}
            {newGroupDescription}
            isCreating={isCreatingGroup}
            error={createGroupError}
        />
    {/if}

    <!-- Join Group Modal -->
    {#if showJoinGroupModal}
        <JoinGroupModal
            on:close={closeJoinGroupModal}
            on:updateCode={(e) => (inviteCode = e.detail)}
            on:preview={handlePreviewGroup}
            on:join={handleJoinGroup}
            {inviteCode}
            {isLoadingPreview}
            {groupPreview}
            {isJoiningGroup}
            error={joinGroupError}
        />
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
