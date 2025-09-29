<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label } from "@nanoteck137/nano-ui";
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

  let { data } = $props();
  const apiClient = getApiClient();

  $effect(() => {
    reset({});
  });

  const { form, errors, enhance, validateForm, reset } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;
          const res = await apiClient.editMedia(data.media.id, {
            coverUrl: formData.coverUrl !== "" ? formData.coverUrl : null,
            bannerUrl: formData.bannerUrl !== "" ? formData.bannerUrl : null,
            logoUrl: formData.logoUrl !== "" ? formData.logoUrl : null,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          toast.success("Successfully updated media images");
          invalidateAll();
        }
      },
    },
  );

  validateForm({ update: true });
</script>

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

  <Button type="submit">Save</Button>
</form>
