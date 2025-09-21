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
  import type { MediaPart } from "$lib/api/types";
  import { isMatch } from "date-fns";

  const Schema = z.object({
    name: z.string().min(1, "Name cannot be empty"),

    releaseDate: z
      .string()
      .trim()
      .refine((val) => val === "" || isMatch(val, "yyyy-MM-dd"), {
        message: "Release date is not in the correct format (YYYY-MM-DD)",
      }),
  });

  export type Props = {
    open: boolean;
    part: MediaPart;
  };

  let { open = $bindable(), part }: Props = $props();
  const apiClient = getApiClient();

  $effect(() => {
    if (open) {
      reset({
        data: {
          name: part.name,
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
          const data = form.data;
          const res = await apiClient.editPart(
            part.mediaId,
            part.index.toString(),
            {
              name: data.name,
            },
          );
          if (!res.success) {
            return handleApiError(res.error);
          }

          open = false;

          toast.success("Successfully updated part");
          invalidateAll();
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[420px] overflow-y-scroll">
    <Dialog.Header>
      <Dialog.Title>Update part</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4 px-[1px]" use:enhance>
      <FormItem>
        <Label for="name">Name</Label>
        <Input id="name" name="name" type="text" bind:value={$form.name} />
        <Errors errors={$errors.name} />
      </FormItem>

      <FormItem>
        <Label for="releaseDate">Release Date (YYYY-MM-DD)</Label>
        <Input
          id="releaseDate"
          name="releaseDate"
          type="text"
          bind:value={$form.releaseDate}
        />
        <Errors errors={$errors.releaseDate} />
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
