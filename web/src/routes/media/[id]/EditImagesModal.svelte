<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import type { ProviderValue } from "$lib/api/types";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import {
    Button,
    Checkbox,
    Dialog,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

  const Schema = z.object({
    coverUrl: z.string().url().optional().or(z.literal("")),
    bannerUrl: z.string().url().optional().or(z.literal("")),
    logoUrl: z.string().url().optional().or(z.literal("")),
  });
  type SchemaTy = z.infer<typeof Schema>;

  export type Props = {
    open: boolean;
    mediaId: string;
  };

  let { open = $bindable(), mediaId }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    reset({ data: { coverUrl: "", bannerUrl: "", logoUrl: "" } });
  });

  async function submit(data: SchemaTy) {
    const res = await apiClient.editMedia(mediaId, {
      coverUrl: data.coverUrl !== "" ? data.coverUrl : null,
      bannerUrl: data.bannerUrl !== "" ? data.bannerUrl : null,
      logoUrl: data.logoUrl !== "" ? data.logoUrl : null,
    });
    if (!res.success) {
      return handleApiError(res.error);
    }
    toast.success("Successfully update media");
    invalidateAll();
  }

  const { form, errors, enhance, validateForm, reset } = superForm(
    defaults({ coverUrl: "", bannerUrl: "", logoUrl: "" }, zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      onUpdate({ form }) {
        if (form.valid) {
          submit(form.data);

          open = false;
          reset({});
        }
      },
    },
  );

  validateForm({ update: true });
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Edit Media Images</Dialog.Title>
      <Dialog.Description>Set the media images</Dialog.Description>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="coverUrl">Cover URL</Label>
        <Input
          id="coverUrl"
          name="coverUrl"
          type="text"
          bind:value={$form.coverUrl}
        />
        <Errors errors={$errors.coverUrl} />
      </FormItem>

      <FormItem>
        <Label for="bannerUrl">Banner URL</Label>
        <Input
          id="bannerUrl"
          name="bannerUrl"
          type="text"
          bind:value={$form.bannerUrl}
        />
        <Errors errors={$errors.bannerUrl} />
      </FormItem>

      <FormItem>
        <Label for="logoUrl">Logo URL</Label>
        <Input
          id="logoUrl"
          name="logoUrl"
          type="text"
          bind:value={$form.logoUrl}
        />
        <Errors errors={$errors.logoUrl} />
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
