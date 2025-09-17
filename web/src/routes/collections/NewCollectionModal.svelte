<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import {
    Button,
    Dialog,
    Input,
    Label,
    Select,
    Separator,
  } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import { CollectionTypeEnum, collectionTypes } from "./types";

  const Schema = z.object({
    type: CollectionTypeEnum.default("unknown"),
    name: z.string().min(1, "Name cannot be empty"),
    coverUrl: z
      .string()
      .url("Cover URL must be valid url")
      .optional()
      .or(z.literal("")),
    bannerUrl: z
      .string()
      .url("Banner URL must be valid url")
      .optional()
      .or(z.literal("")),
    logoUrl: z
      .string()
      .url("Logo URL must be valid url")
      .optional()
      .or(z.literal("")),
  });
  type SchemaTy = z.infer<typeof Schema>;

  export type Props = {
    open: boolean;
  };

  let { open = $bindable() }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({});
    }
  });

  async function submit(data: SchemaTy) {
    const res = await apiClient.createCollection({
      name: data.name,
      collectionType: data.type,
      coverUrl: data.coverUrl ?? "",
      bannerUrl: data.bannerUrl ?? "",
      logoUrl: data.logoUrl ?? "",
    });
    if (!res.success) {
      return handleApiError(res.error);
    }

    toast.success("Successfully updated media");
    invalidateAll();
  }

  const { form, errors, enhance, reset } = superForm(defaults(zod(Schema)), {
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
  });
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Create new collection</Dialog.Title>
      <!-- <Dialog.Description>Set the media images</Dialog.Description> -->
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="type">Type</Label>

        <Select.Root type="single" bind:value={$form.type}>
          <Select.Trigger>
            {collectionTypes.find((f) => f.value === $form.type)?.label ??
              "Select type"}
          </Select.Trigger>
          <Select.Content>
            {#each collectionTypes as type (type.value)}
              <Select.Item value={type.value} label={type.label} />
            {/each}
          </Select.Content>
        </Select.Root>

        <Errors errors={$errors.type} />
      </FormItem>

      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </FormItem>

      <Separator />

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

        <Button type="submit">Create</Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
