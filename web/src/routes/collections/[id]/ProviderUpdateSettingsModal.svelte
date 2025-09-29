<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { ProviderValue } from "$lib/api/types";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Checkbox, Dialog, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

  const Schema = z.object({
    replaceImages: z.boolean(),
  });

  export type Props = {
    open: boolean;
    collectionId: string;
    provider: ProviderValue;
  };

  let { open = $bindable(), collectionId, provider }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({ data: { replaceImages: false } });
    }
  });

  const { form, errors, enhance, reset } = superForm(defaults(zod(Schema)), {
    SPA: true,
    validators: zod(Schema),
    resetForm: true,
    async onUpdate({ form }) {
      if (form.valid) {
        const formData = form.data;
        const res = await apiClient.providerUpdateCollection(
          provider.name,
          collectionId,
          formData,
        );
        if (!res.success) {
          return handleApiError(res.error);
        }

        open = false;

        toast.success("Successfully update media");
        invalidateAll();
      }
    },
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Update media: {provider.displayName}</Dialog.Title>
      <Dialog.Description>
        Use '{provider.displayName}' to update the media
      </Dialog.Description>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <div class="flex items-center gap-2">
          <Checkbox
            id="replaceImages"
            name="replaceImages"
            bind:checked={$form.replaceImages}
          />
          <Label for="replaceImages">Replace Images</Label>
        </div>
        <Errors errors={$errors.replaceImages} />
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

        <Button type="submit">Save</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
