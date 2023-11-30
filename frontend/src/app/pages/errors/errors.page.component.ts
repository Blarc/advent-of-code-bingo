import {NgIf} from '@angular/common';
import {Component} from '@angular/core';
import {Router, RouterLink} from '@angular/router';

@Component({
    standalone: true,
    selector: 'app-errors',
    styleUrls: ['errors.page.component.scss'],
    imports: [NgIf, RouterLink],
    templateUrl: 'errors.page.component.html'
})
export class ErrorsPageComponent {
    public error = false;

    constructor(private router: Router) {
        this.error = this.router.getCurrentNavigation()?.extras.state?.['error'];
    }
}
