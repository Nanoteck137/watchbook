<script lang="ts">
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/modals";
  import {
    Button,
    Dialog,
    Input,
    Label,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import { UndoDot } from "lucide-svelte";
  import type { MediaItem, SearchResult } from "./AddMediaItem.svelte";
  import Errors from "$lib/components/Errors.svelte";
  import Spacer from "$lib/components/Spacer.svelte";
  import { z } from "zod";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { zod } from "sveltekit-superforms/adapters";

  const Schema = z.object({
    items: z.array(
      z.object({
        name: z.string().min(1),
        position: z.number().min(0),
      }),
    ),
  });

  export type Props = {
    open: boolean;
    items: SearchResult[];
  };

  let {
    open = $bindable(),
    onResult,

    items,
  }: Props & Modal<MediaItem[]> = $props();

  $effect(() => {
    reset({
      data: {
        items: items.map((i) => ({
          name: i.name,
          searchSlug: "",
          position: 0,
        })),
      },
    });
  });

  const { form, errors, enhance, validateForm, reset } = superForm(
    defaults(zod(Schema)),
    {
      SPA: true,
      validators: zod(Schema),
      resetForm: true,
      dataType: "json",
      onUpdate({ form }) {
        if (form.valid) {
          onResult(
            form.data.items.map((item, i) => ({
              id: items[i].id,
              name: item.name,
              position: item.position,
            })),
          );
          open = false;
          reset({});
        }
      },
    },
  );

  validateForm({ update: true });
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[420px] overflow-y-scroll">
    <Dialog.Header>
      <Dialog.Title>Add media to the season</Dialog.Title>
    </Dialog.Header>
    <form class="flex flex-col gap-4" method="POST" use:enhance>
      {#each $form.items as _item, i}
        <div class="flex flex-col gap-4 border-t pt-2">
          <h2 class="text-base">{items[i].name}</h2>

          <FormItem>
            <Label for="position">Position</Label>
            <Input
              id="position"
              name="position"
              type="number"
              bind:value={$form.items[i].position}
            />
            <Errors errors={$errors.items?.[i]?.position} />
          </FormItem>
        </div>
      {/each}

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
