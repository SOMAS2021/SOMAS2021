import { Intent, Position, Toaster } from "@blueprintjs/core";

const AppToaster = Toaster.create({
  className: "recipe-toaster",
  position: Position.TOP,
});

export const showToast = (message: string, intent: Intent) => {
  // create toasts in response to interactions.
  // in most cases, it's enough to simply create and forget (thanks to timeout).
  AppToaster.show({ message: message, intent: intent });
};
