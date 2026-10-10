# done

A terminal task finisher.

## Commands

```text
done                      # open the task list
done add "Write README"   # add a task
done ls                   # list open tasks
done check <task id>      # mark a task as completed
done clear                # remove completed tasks
```

In the task list, use the arrow keys to move, space to toggle a task, `a` to add one,
and `q` to quit. The footer shows today's completion progress; quitting prints a short
session summary.

Tasks are stored as JSON at `done/tasks.json` under the user's config directory.

## Demo

With [VHS](https://github.com/charmbracelet/vhs) installed, run `vhs demo.tape` from the project directory to make `demo.gif`.
