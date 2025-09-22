<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { ProviderValue } from "$lib/api/types";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import Spinner from "$lib/components/Spinner.svelte";
  import { Button, Checkbox, Dialog, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

  const Schema = z.object({
    replaceImages: z.boolean(),
    overrideParts: z.boolean().default(true),
    setRelease: z.boolean().default(false),
  });

  export type Props = {
    open: boolean;
    mediaId: string;
    provider: ProviderValue;
  };

  let { open = $bindable(), mediaId, provider }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset();
    }
  });

  const { form, errors, enhance, validateForm, reset, submitting } = superForm(
    defaults({}, zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const data = form.data;
          const res = await apiClient.providerUpdateMedia(
            provider.name,
            mediaId,
            data,
          );
          if (!res.success) {
            return handleApiError(res.error);
          }

          toast.success("Successfully update media");
          invalidateAll();

          open = false;
        }
      },
    },
  );

  validateForm({ update: true });
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

      <FormItem>
        <div class="flex items-center gap-2">
          <Checkbox
            id="overrideParts"
            name="overrideParts"
            bind:checked={$form.overrideParts}
          />
          <Label for="overrideParts">Override Parts</Label>
        </div>
        <Errors errors={$errors.overrideParts} />
      </FormItem>

      <FormItem>
        <div class="flex items-center gap-2">
          <Checkbox
            id="setRelease"
            name="setRelease"
            bind:checked={$form.setRelease}
          />
          <Label for="setRelease">Set Release</Label>
        </div>
        <Errors errors={$errors.setRelease} />
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
