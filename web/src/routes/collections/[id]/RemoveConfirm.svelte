<script lang="ts">
  import type { Modal } from "$lib/components/modals";
  import { Button, Dialog } from "@nanoteck137/nano-ui";

  export type Props = {
    open: boolean;
  };

  let { open = $bindable(), onResult }: Props & Modal<void> = $props();
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Remove Collection Item?</Dialog.Title>
      <Dialog.Description>
        Are you sure you want to remove this media item? This action cannot be
        undone.
      </Dialog.Description>
    </Dialog.Header>

    <Dialog.Footer class="gap-2 sm:gap-0">
      <Button
        variant="outline"
        onclick={() => {
          open = false;
        }}
      >
        Cancel
      </Button>

      <Button
        variant="destructive"
        onclick={() => {
          open = false;
          onResult();
        }}
      >
        Remove item
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
