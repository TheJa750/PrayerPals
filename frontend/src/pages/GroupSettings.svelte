<script>
    import { onMount } from "svelte";
    import { apiRequest } from "../lib/api";
    import { REFRESH_ERROR_MESSAGE } from "../lib/api";
    import ChangeInviteCodeModal from "../components/ChangeInviteCodeModal.svelte";
    import { validateInviteCode } from "../lib/validation";
    import ChangeGroupRulesModal from "../components/ChangeGroupRulesModal.svelte";
    import { marked } from "marked";
    import DOMPurify from "dompurify";

    export let navigate;
    export let groupId;

    // Group data variables
    let group = null;
    let isLoadingGroup = false;
    let loadGroupError = "";
    let groupRulesHTML = null;

    // Modal state variables
    let showChangeCodeModal = false;
    let showChangeRulesModal = false;

    // Invite code variables
    let newInviteCode = "";
    let inviteCodeError = "";
    let isChangingCode = false;
    let newInviteCodeValidation = { isValid: true, errors: [] };

    // Group rules variables
    let groupRulesMD = "";
    let isSavingRules = false;
    let rulesError = "";

    async function loadGroupData() {
        try {
            isLoadingGroup = true;
            loadGroupError = "";

            const groupData = await apiRequest(`/groups/${groupId}`, "GET");
            group = groupData;

            groupRulesMD = groupData.rules_info || "";
            groupRulesHTML = marked(groupRulesMD) || "";
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

    // Modal functions
    function openChangeCodeModal() {
        showChangeCodeModal = true;
    }

    function closeChangeCodeModal() {
        showChangeCodeModal = false;
        newInviteCode = "";
        inviteCodeError = "";
        newInviteCodeValidation = { isValid: true, errors: [] };
    }

    function openChangeRulesModal() {
        showChangeRulesModal = true;
    }

    function closeChangeRulesModal() {
        showChangeRulesModal = false;
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

    function handleChangeRulesUpdate(event) {
        groupRulesMD = event.detail;
    }

    async function handleSaveRules(event) {
        event.preventDefault();
        isSavingRules = true;
        rulesError = "";

        try {
            await apiRequest(`/groups/${groupId}/rules`, "PUT", {
                rules: groupRulesMD,
            });
            closeChangeRulesModal();
            await loadGroupData(); // Refresh group data
        } catch (error) {
            console.error("Error saving group rules:", error);
            if (error.message === REFRESH_ERROR_MESSAGE) {
                navigate("login");
                return;
            }
            rulesError =
                error.message ||
                "Failed to save group rules. Please try again.";
        } finally {
            isSavingRules = false;
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
        <div class="settings-card">
            <h2>Group Rules</h2>
            <p>Update your group's rules. Use Markdown for formatting.</p>
            <div class="group-rules rules-content">
                {#if group.rules_info}
                    <p>{@html DOMPurify.sanitize(groupRulesHTML)}</p>
                {:else}
                    <p>No rules set for this group.</p>
                {/if}
            </div>
            <div class="card-action-row justify-right">
                <button class="action-button" on:click={openChangeRulesModal}
                    >Edit Rules</button
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

{#if showChangeRulesModal}
    <ChangeGroupRulesModal
        groupRules={groupRulesMD}
        error={rulesError}
        isSaving={isSavingRules}
        on:close={closeChangeRulesModal}
        on:updateRules={handleChangeRulesUpdate}
        on:submit={handleSaveRules}
    />
{/if}
