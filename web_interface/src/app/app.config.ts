import { ApplicationConfig, provideBrowserGlobalErrorListeners, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { LucideAngularModule, Calendar, Clock, CircleAlert, Plus, EllipsisVertical, X, CircleCheck, ListTodo, CirclePlay, Moon, Sun } from 'lucide-angular';


import { routes } from './app.routes';
import { provideClientHydration, withEventReplay } from '@angular/platform-browser';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideRouter(routes), provideClientHydration(withEventReplay()),
    provideHttpClient(withFetch()),
    importProvidersFrom(LucideAngularModule.pick({ Calendar, Clock, CircleAlert, Plus, EllipsisVertical, X, CircleCheck, ListTodo, CirclePlay, Moon, Sun }))
  ]
};
