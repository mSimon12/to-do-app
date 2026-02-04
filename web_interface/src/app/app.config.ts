import { ApplicationConfig, provideBrowserGlobalErrorListeners, importProvidersFrom } from '@angular/core';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withFetch } from '@angular/common/http';
import { LucideAngularModule, Calendar, Clock, AlertCircle, Plus, MoreVertical, X, CheckCircle, ListTodo, PlayCircle } from 'lucide-angular';


import { routes } from './app.routes';
import { provideClientHydration, withEventReplay } from '@angular/platform-browser';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideHttpClient(withFetch()),
    importProvidersFrom(LucideAngularModule.pick({ Calendar, Clock, AlertCircle, Plus, MoreVertical, X, CheckCircle, ListTodo, PlayCircle })),
    provideRouter(routes), provideClientHydration(withEventReplay())
  ]
};
