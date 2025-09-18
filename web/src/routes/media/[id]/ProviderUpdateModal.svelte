<script lang="ts">
  import { Button, Dialog } from "@nanoteck137/nano-ui";
  import type { ProviderValue } from "$lib/api/types";
  import ProviderUpdate from "./ProviderUpdateButton.svelte";

  export type Props = {
    open: boolean;
    mediaId: string;
    providers: ProviderValue[];
  };

  let { open = $bindable(), mediaId, providers }: Props = $props();
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Provider Update</Dialog.Title>
      <Dialog.Description>
        Choose the provider used to update the media
      </Dialog.Description>
    </Dialog.Header>

    {#each providers as provider}
      <ProviderUpdate {mediaId} {provider} />
    {/each}

    <Dialog.Footer class="gap-2 sm:gap-0">
      <Button
        variant="outline"
        onclick={() => {
          open = false;
        }}
      >
        Close
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
