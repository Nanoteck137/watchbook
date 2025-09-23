<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Input, Label, Select } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { Collection } from "$lib/api/types";
  import {
    type CollectionType,
    CollectionTypeEnum,
    collectionTypes,
  } from "$lib/api-types";

  const Schema = z.object({
    type: CollectionTypeEnum.default("unknown"),
    name: z.string().min(1, "Name cannot be empty"),
  });

  let { data } = $props();
  const apiClient = getApiClient();

  $effect(() => {
    reset({
      data: {
        type: data.collection.type as CollectionType,
        name: data.collection.name,
      },
    });
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(
      {
        type: data.collection.type as CollectionType,
        name: data.collection.name,
      },
      zod(Schema),
    ),
    {
      id: "edit",
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          const formData = form.data;
          const res = await apiClient.editCollection(data.collection.id, {
            type: formData.type,
            name: formData.name,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          toast.success("Successfully updated collection");
          invalidateAll();
        }
      },
    },
  );
</script>

<form class="flex flex-col gap-4" use:enhance>
  <FormItem>
    <Label for="type">Type</Label>

    <Select.Root type="single" bind:value={$form.type} allowDeselect={false}>
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

  <Button type="submit" disabled={$submitting}>
    Update
    {#if $submitting}
      <Spinner />
    {/if}
  </Button>
</form>
