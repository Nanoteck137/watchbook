<script lang="ts">
  import FormItem from "$lib/components/FormItem.svelte";
  import { Input, Label } from "@nanoteck137/nano-ui";
</script>

<p>New Entry</p>

<FormItem>
  <Label for="title">Title</Label>
  <Input id="title" name="title" type="text" />
  <!-- <Errors errors={$errors.username} /> -->
</FormItem>
