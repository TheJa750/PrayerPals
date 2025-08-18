<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import { REFRESH_ERROR_MESSAGE } from "../lib/api";
    import ChangeInviteCodeModal from "../components/ChangeInviteCodeModal.svelte";
    import { validateInviteCode } from "../lib/validation";

    export let navigate;
    export let groupId;

    // Group data variables
    let group = null;
    let isLoadingGroup = false;
    let loadGroupError = "";

    // Modal state variables
    let showChangeCodeModal = false;

    // Invite code variables
    let newInviteCode = "";
    let inviteCodeError = "";
    let isChangingCode = false;
    let newInviteCodeValidation = { isValid: true, errors: [] };

    async function loadGroupData() {
        try {
            isLoadingGroup = true;
            loadGroupError = "";

            const groupData = await apiRequest(`/groups/${groupId}`, "GET");
            group = groupData;
        } catch (error) {
            console.error("Error loading group data:", error);

            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }

            loadGroupError =
                "Failed to load group data. Please try again later.";
        } finally {
            isLoadingGroup = false;
        }
    }

    onMount(() => {
        if (groupId) {
            loadGroupData();
        } else {
            navigate("user");
        }
    });

    function openChangeCodeModal() {
        showChangeCodeModal = true;
    }

    function closeChangeCodeModal() {
        showChangeCodeModal = false;
        newInviteCode = "";
        inviteCodeError = "";
        newInviteCodeValidation = { isValid: true, errors: [] };
    }

    function handleInviteCodeUpdate(event) {
        newInviteCode = event.detail;
        newInviteCodeValidation = validateInviteCode(newInviteCode);
        inviteCodeError = newInviteCodeValidation.errors.join(" ");
    }

    async function handleChangeCode(event) {
        event.preventDefault();
        if (!newInviteCodeValidation.isValid) {
            inviteCodeError = newInviteCodeValidation.errors.join(", ");
            return;
        }

        isChangingCode = true;
        try {
            await apiRequest(`/groups/${groupId}/invite-code`, "PUT", {
                invite_code: newInviteCode,
            });
            group.invite_code = newInviteCode; // Update local state
            closeChangeCodeModal();
            await loadGroupData(); // Refresh group data
        } catch (error) {
            console.error("Error changing invite code:", error);
            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }
            inviteCodeError =
                error.message ||
                "Failed to change invite code. Please try again.";
        } finally {
            isChangingCode = false;
        }
    }
</script>

<div class="group-settings">
    <div class="header-row">
        <button class="back-button" on:click={() => navigate("group", groupId)}>
            ←
        </button>
        <div class="header-content">
            <h2>Group Settings</h2>
            <p>Manage your group settings here.</p>
        </div>
    </div>
    {#if isLoadingGroup}
        <p>Loading group settings...</p>
    {:else if loadGroupError}
        <p class="error">{loadGroupError}</p>
    {:else if group}
        <div class="settings-card">
            <h2>Change Invitation Code</h2>
            <p>
                Update your group's invitation code. (Max length: 6 characters)
            </p>
            <p>
                Note: The final invitation code will be 9 characters long. We
                will randomly generate at least 3 characters to prevent
                duplicate codes.
            </p>
            <div class="card-action-row">
                <h3>Current Code: {group.invite_code}</h3>
                <button class="action-button" on:click={openChangeCodeModal}
                    >Change Code</button
                >
            </div>
        </div>
    {/if}
</div>

{#if showChangeCodeModal}
    <ChangeInviteCodeModal
        inviteCode={newInviteCode}
        error={inviteCodeError}
        {isChangingCode}
        on:close={closeChangeCodeModal}
        on:updateCode={handleInviteCodeUpdate}
        on:submit={handleChangeCode}
    />
{/if}
