<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label, Select } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import {
    CollectionTypeEnum,
    collectionTypes,
    type CollectionType,
  } from "../types";
  import Spinner from "$lib/components/Spinner.svelte";
  import type { Collection } from "$lib/api/types";

  const Schema = z.object({
    type: CollectionTypeEnum.default("unknown"),
    name: z.string().min(1, "Name cannot be empty"),
  });
  type SchemaTy = z.infer<typeof Schema>;

  export type Props = {
    open: boolean;
    collection: Collection;
  };

  let { open = $bindable(), collection }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({
        data: {
          type: collection.collectionType as CollectionType,
          name: collection.name,
        },
      });
    }
  });

  async function submit(data: SchemaTy) {
    const res = await apiClient.editCollection(collection.id, {
      collectionType: data.type,
      name: data.name,
    });
    if (!res.success) {
      return handleApiError(res.error);
    }

    open = false;
    toast.success("Successfully updated collection");
    invalidateAll();
  }

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(
      {
        type: collection.collectionType as CollectionType,
        name: collection.name,
      },
      zod(Schema),
    ),
    {
      SPA: true,
      validators: zod(Schema),
      dataType: "json",
      resetForm: true,
      async onUpdate({ form }) {
        if (form.valid) {
          await submit(form.data);
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Edit collection</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="type">Type</Label>

        <Select.Root
          type="single"
          bind:value={$form.type}
          allowDeselect={false}
        >
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
