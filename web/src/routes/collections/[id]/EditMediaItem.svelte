<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/modals";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

  const schema = z.object({
    name: z.string().min(1, "String must contain at least 1 character(s)"),
    searchSlug: z.string(),
    order: z.number().min(0, "Order cannot be negative"),
  });

  type Result = {
    name: string;
    searchSlug: string;
    order: number;
  };

  export type Props = {
    open: boolean;
    name: string;
    searchSlug: string;
    order: number;
  };

  let {
    open = $bindable(),

    name,
    searchSlug,
    order,
    onResult,
  }: Props & Modal<Result> = $props();

  $effect(() => {
    reset({ data: { name, searchSlug, order } });
  });

  const { form, errors, enhance, validateForm, reset } = superForm(
    defaults({ name, searchSlug, order }, zod(schema)),
    {
      SPA: true,
      validators: zod(schema),
      resetForm: true,
      onUpdate({ form }) {
        console.log(form);
        if (form.valid) {
          onResult(form.data);
          console.log("Valid form", form.data);
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
      <Dialog.Title>Edit collection media item</Dialog.Title>
      <Dialog.Description>
        Update the details of this collection media item. Make changes and save
        when you’re done.
      </Dialog.Description>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </FormItem>

      <div class="flex flex-col-reverse gap-4 sm:flex-row">
        <FormItem class="w-full sm:w-32">
          <Label for="order">Order</Label>
          <Input
            id="order"
            name="order"
            type="number"
            bind:value={$form.order}
          />
          <Errors errors={$errors.order} />
        </FormItem>

        <FormItem class="w-full">
          <Label for="searchSlug">Search Slug</Label>
          <Input
            id="searchSlug"
            name="searchSlug"
            type="text"
            bind:value={$form.searchSlug}
          />
          <Errors errors={$errors.searchSlug} />
        </FormItem>
      </div>

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
