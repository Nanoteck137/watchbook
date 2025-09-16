<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/modals";
  import {
    Button,
    Checkbox,
    Dialog,
    Input,
    Label,
  } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

  const schema = z.object({
    replaceImages: z.boolean(),
  });

  type Result = {
    replaceImages: boolean;
  };

  export type Props = {
    open: boolean;
    providerDisplayName: string;
  };

  let {
    open = $bindable(),

    providerDisplayName,

    onResult,
  }: Props & Modal<Result> = $props();

  $effect(() => {
    reset({ data: { replaceImages: false } });
  });

  const { form, errors, enhance, validateForm, reset } = superForm(
    defaults({ replaceImages: false }, zod(schema)),
    {
      SPA: true,
      validators: zod(schema),
      resetForm: true,
      onUpdate({ form }) {
        console.log(form);
        if (form.valid) {
          onResult(form.data);
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
      <Dialog.Title>Update media: {providerDisplayName}</Dialog.Title>
      <Dialog.Description>
        Use '{providerDisplayName}' to update the media
      </Dialog.Description>
    </Dialog.Header>

    <form class="flex flex-col gap-4" use:enhance>
      <FormItem>
        <div class="flex items-center gap-2">
          <Checkbox
            id="replaceImages"
            name="replaceImages"
            bind:checked={$form.replaceImages}
          />
          <Label for="replaceImages">Replace Images</Label>
        </div>
        <!-- <Input type="text" /> -->
        <Errors errors={$errors.replaceImages} />
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
