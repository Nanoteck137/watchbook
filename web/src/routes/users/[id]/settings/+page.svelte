<script lang="ts">
  import { Button, Card } from "@nanoteck137/nano-ui";
  import ChangePassword from "./ChangePassword.svelte";
  import ChangeDisplayName from "./ChangeDisplayName.svelte";
  import ApiToken from "./ApiToken.svelte";
  import NewApiTokenModal from "./NewApiTokenModal.svelte";
  import { Plus } from "lucide-svelte";

  let { data } = $props();

  let openNewApiTokenModal = $state(false);
</script>

<Card.Root>
  <ChangeDisplayName />
  <ChangePassword />

  <div class="flex flex-col items-center gap-4 border-b p-6">
    <h2 class="text-bold text-center text-xl">
      API Tokens
      <Button
        variant="ghost"
        size="icon"
        onclick={() => {
          openNewApiTokenModal = true;
        }}
      >
        <Plus />
      </Button>
    </h2>

    <div class="flex w-full flex-col sm:max-w-[460px]">
      {#each data.tokens as token}
        <ApiToken {token} />
      {/each}
    </div>
  </div>
</Card.Root>

<NewApiTokenModal bind:open={openNewApiTokenModal} />
