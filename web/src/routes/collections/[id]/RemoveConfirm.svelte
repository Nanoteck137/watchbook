<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/modals";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

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
