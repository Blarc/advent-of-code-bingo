import {AsyncPipe, NgIf} from '@angular/common';
import {Component, OnInit, inject} from '@angular/core';
import {RouterOutlet} from '@angular/router';

import {Observable} from 'rxjs';

import {UserDto} from '../../../api/model/userDto.model';
import {AuthService} from '../../../services/auth.service';

@Component({
    standalone: true,
    selector: 'app-header',
    templateUrl: 'header.component.html',
    imports: [RouterOutlet, NgIf, AsyncPipe],
    styleUrls: ['header.component.scss']
})
export class HeaderComponent implements OnInit {
    private authService = inject(AuthService);

    public user?: Observable<UserDto | null>;

    ngOnInit() {
        this.user = this.authService.userSubject$;
    }
}
