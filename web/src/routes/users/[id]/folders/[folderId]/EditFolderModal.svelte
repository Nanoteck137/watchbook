<script lang="ts">
  import { goto, invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { Folder } from "$lib/api/types";
  import ConfirmBox from "$lib/components/ConfirmBox.svelte";
  import { Trash } from "lucide-svelte";

  const Schema = z.object({
    name: z.string().min(1),
  });

  export type Props = {
    open: boolean;
    folder: Folder;
  };

  let { open = $bindable(), folder }: Props = $props();
  const apiClient = getApiClient();

  let openDeleteConfirmModal = $state(false);

  $effect(() => {
    if (open) {
      reset({
        data: {
          name: folder.name,
        },
      });
    }
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;

          const res = await apiClient.editFolder(folder.id, {
            name: formData.name,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully updated folder");
          invalidateAll();
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[420px] overflow-y-scroll">
    <Dialog.Header>
      <Dialog.Title>Edit folder</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4 px-[1px]" use:enhance>
      <FormItem>
        <Label for="name">Name</Label>
        <div class="flex items-center gap-2">
          <Input id="name" name="name" type="text" bind:value={$form.name} />
          <Button
            variant="destructive"
            onclick={() => {
              openDeleteConfirmModal = true;
            }}
          >
            <Trash />
          </Button>
        </div>
        <Errors errors={$errors.name} />
      </FormItem>

      <Dialog.Footer class="gap-2 sm:gap-0">
        <Button
          variant="outline"
          onclick={() => {
            open = false;
            reset();
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Update
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>

<ConfirmBox
  bind:open={openDeleteConfirmModal}
  title="Delete this folder?"
  description="Are you sure you want to delete this folder? This action cannot be undone."
  confirmText="Delete Folder"
  onResult={async () => {
    const res = await apiClient.deleteFolder(folder.id);
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully deleted folder");
    goto(`/users/${folder.userId}/folders`, { invalidateAll: true });
  }}
/>
