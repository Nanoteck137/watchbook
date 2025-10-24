<script lang="ts">
  import { invalidateAll } from "$app/navigation";
  import { getApiClient, handleApiError } from "$lib";
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";
  import Spinner from "$lib/components/Spinner.svelte";

  const Schema = z.object({
    num: z.number().min(0),
    name: z.string().min(1),
    searchSlug: z.string(),
  });

  export type Props = {
    open: boolean;
    showId: string;
  };

  let { open = $bindable(), showId }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({});
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
          const res = await apiClient.addShowSeason(showId, {
            num: formData.num,
            name: formData.name,
            searchSlug: formData.searchSlug,
          });
          if (!res.success) {
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully created show season");
          invalidateAll();
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Add new Season</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <Label for="num">Season Number</Label>
        <Input id="num" name="num" type="number" bind:value={$form.num} />
        <Errors errors={$errors.num} />
      </FormItem>

      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </FormItem>

      <FormItem>
        <Label for="searchSlug">Search Slug (defaults to name)</Label>
        <Input
          id="searchSlug"
          name="searchSlug"
          type="text"
          bind:value={$form.searchSlug}
        />
        <Errors errors={$errors.searchSlug} />
      </FormItem>

      <Dialog.Footer class="gap-2 sm:gap-0">
        <Button
          variant="outline"
          onclick={() => {
            open = false;
          }}
        >
          Close
        </Button>

        <Button type="submit" disabled={$submitting}>
          Create
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
