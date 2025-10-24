<script lang="ts">
  import Errors from "$lib/components/Errors.svelte";
  import FormItem from "$lib/components/FormItem.svelte";
  import type { Modal } from "$lib/components/modals";
  import Spinner from "$lib/components/Spinner.svelte";
  import { Button, Dialog, Input, Label } from "@nanoteck137/nano-ui";
  import { zod } from "sveltekit-superforms/adapters";
  import { defaults, superForm } from "sveltekit-superforms/client";
  import { z } from "zod";

  const schema = z.object({
    position: z.number().min(0),
  });

  type Result = {
    position: number;
  };

  export type Props = {
    open: boolean;
    position: number;
  };

  let {
    open = $bindable(),

    position,
    onResult,
  }: Props & Modal<Result> = $props();

  $effect(() => {
    reset({ data: { position } });
  });

  const { form, errors, enhance, reset, submitting } = superForm(
    defaults(zod(schema)),
    {
      SPA: true,
      validators: zod(schema),
      dataType: "json",
      resetForm: true,
      onUpdate({ form }) {
        if (form.valid) {
          onResult(form.data);
          open = false;
        }
      },
    },
  );
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-h-[420px] overflow-y-scroll">
    <Dialog.Header>
      <Dialog.Title>Edit media item</Dialog.Title>
    </Dialog.Header>

    <form class="flex flex-col gap-4 px-[1px]" use:enhance>
      <FormItem>
        <Label for="position">Position</Label>
        <Input
          id="position"
          name="position"
          type="number"
          bind:value={$form.position}
        />
        <Errors errors={$errors.position} />
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
          Save
          {#if $submitting}
            <Spinner />
          {/if}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
