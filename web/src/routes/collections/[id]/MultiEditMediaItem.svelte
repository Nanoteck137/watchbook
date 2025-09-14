<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import FormItem from "$lib/components/FormItem.svelte";
  import Image from "$lib/components/Image.svelte";
  import type { Modal } from "$lib/components/modals";
  import { cn } from "$lib/utils";
  import {
    Button,
    Dialog,
    Input,
    Label,
    ScrollArea,
  } from "@nanoteck137/nano-ui";
  import { Check, UndoDot } from "lucide-svelte";
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
        searchSlug: z.string(),
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
              searchSlug: item.searchSlug,
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
  <Dialog.Content>
    <Dialog.Header>
      <Dialog.Title>Add media to the collection</Dialog.Title>
    </Dialog.Header>
    <form class="flex flex-col gap-2" method="POST" use:enhance>
      <ScrollArea class="h-60">
        <div class="px-[1px]">
          {#each $form.items as item, i}
            <div class="flex flex-col gap-4 border-t">
              <div>
                <Spacer />

                <Label for="name">{items[i].name}</Label>
              </div>

              <FormItem>
                <Label for="name">Name</Label>
                <Input
                  id="name"
                  name="name"
                  type="text"
                  bind:value={$form.items[i].name}
                />
                <Errors errors={$errors.items?.[i]?.name} />
              </FormItem>

              <div class="flex flex-col-reverse gap-4 sm:flex-row">
                <FormItem class="w-full sm:w-32">
                  <Label class="flex h-5 items-center" for="position"
                    >Position</Label
                  >
                  <Input
                    id="position"
                    name="position"
                    type="number"
                    bind:value={$form.items[i].position}
                  />
                  <Errors errors={$errors.items?.[i]?.position} />
                </FormItem>

                <FormItem class="w-full">
                  <Label class="flex items-center gap-2" for="searchSlug"
                    >Search Slug

                    <button
                      type="button"
                      onclick={() => {
                        const name = $form.items[i].name;
                        $form.items[i].searchSlug = name;
                      }}
                    >
                      <UndoDot class="h-5 w-5" />
                    </button>
                  </Label>
                  <div class="relative">
                    <Input
                      id="searchSlug"
                      name="searchSlug"
                      type="text"
                      bind:value={$form.items[i].searchSlug}
                    />
                    <button
                      class="absolute bottom-1/2 right-2 translate-y-1/2"
                    >
                      <UndoDot class="h-5 w-5" />
                    </button>
                  </div>
                  <Errors errors={$errors.items?.[i]?.searchSlug} />
                </FormItem>
              </div>
            </div>
          {/each}
        </div>
      </ScrollArea>

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
